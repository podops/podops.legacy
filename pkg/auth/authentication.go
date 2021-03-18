package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fupas/commons/pkg/util"
	"github.com/podops/podops/apiv1"
)

// ResetAccountChallenge creates a new confirmation token and resets the timer
func ResetAccountChallenge(ctx context.Context, account *Account) (*Account, error) {
	// FIXME add a mutex
	// FIXME this is crude!
	token, _ := util.ShortUUID()
	account.Expires = util.IncT(util.Timestamp(), DefaultAuthenticationExpiration)
	account.Ext1 = token
	account.Status = AccountUnconfirmed

	if err := UpdateAccount(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

// ResetAuthToken creates a new authorization token and resets the timer
func ResetAuthToken(ctx context.Context, account *Account) (*Account, error) {
	// FIXME add a mutex
	// FIXME this is crude!
	token, _ := util.ShortUUID()
	account.Expires = util.IncT(util.Timestamp(), DefaultAuthenticationExpiration)
	account.Ext2 = token
	account.Status = AccountLoggedOut

	if err := UpdateAccount(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func Logout(ctx context.Context, account *Account) error {
	// FIXME add a mutex
	auth, err := LookupAuthorization(ctx, account.Realm, account.ClientID)
	if err != nil {
		return err
	}
	if auth != nil {
		auth.Revoked = true
		err = UpdateAuthorization(ctx, auth)
		if err != nil {
			return err
		}
	}

	account.Status = AccountLoggedOut
	return UpdateAccount(ctx, account)
}

// SendAccountChallenge sends a notification to the user promting to confirm the account
func SendAccountChallenge(ctx context.Context, account *Account) error {
	url := fmt.Sprintf("%s/login/%s", apiv1.DefaultAPIEndpoint, account.Ext1)
	fmt.Println("account confirm: " + url)

	// FIXME this is not done!
	return nil
	//return fmt.Errorf("SendLoginChallenge: not implemented")
}

// SendAuthToken sends a notification to the user with the current authentication token
func SendAuthToken(ctx context.Context, account *Account) error {
	// FIXME this is not done, just a crude implementation
	fmt.Println("auth token=" + account.Ext2)
	return nil

	//return fmt.Errorf("SendAuthToken: not implemented")
}

// ConfirmLoginChallenge confirms the account
func ConfirmLoginChallenge(ctx context.Context, token string) (*Account, int, error) {
	// FIXME add a mutex
	if token == "" {
		return nil, http.StatusUnauthorized, fmt.Errorf("Invalid token")
	}

	account, err := FindAccountByToken(ctx, token)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if account == nil {
		return nil, http.StatusNotFound, nil
	}
	now := util.Timestamp()
	if account.Expires < now {
		return account, http.StatusForbidden, nil
	}

	account.Confirmed = now
	account.Status = AccountLoggedOut
	account.Ext1 = ""

	err = UpdateAccount(ctx, account)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return account, http.StatusNoContent, nil
}

// exchangeToken confirms the temporary auth token and creates the permanent one
func exchangeToken(ctx context.Context, req *AuthorizationRequest, loginFrom string) (*Authorization, int, error) {
	// FIXME add a mutex
	var auth *Authorization

	account, err := LookupAccount(ctx, req.Realm, req.UserID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if account == nil {
		return nil, http.StatusNotFound, nil
	}
	now := util.Timestamp()
	if account.Expires < now || account.Ext2 != req.Token {
		return nil, http.StatusUnauthorized, nil
	}

	// all OK, create or update the authorization
	auth, err = LookupAuthorization(ctx, req.Realm, req.ClientID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if auth == nil {
		// FIXME this is hardcoded for podops, make it configurable
		auth = &Authorization{
			ClientID:  req.ClientID,
			Realm:     req.Realm,
			Name:      "DEPRECATED",
			TokenType: DefaultTokenType,
			UserID:    req.UserID,
			Scope:     DefaultAPIScope,
			AuthType:  "DEPRECATED",
			Revoked:   false,
			Created:   now,
		}
	}
	token, _ := util.UUID()

	auth.Token = token
	auth.Expires = now + (DefaultAuthorizationExpiration * 86400)
	auth.Updated = now

	err = CreateAuthorization(ctx, auth)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// update the account
	account.Status = AccountActive
	account.LastLogin = now
	account.LoginCount = account.LoginCount + 1
	account.LoginFrom = loginFrom
	account.Ext1 = ""
	account.Ext2 = ""
	account.Expires = 0

	err = UpdateAccount(ctx, account)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return auth, http.StatusOK, nil
}
