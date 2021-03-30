package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fupas/commons/pkg/env"
	"github.com/fupas/commons/pkg/util"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/platform"
)

const (
	DefaultScope     = "production:read,production:write,production:build,resource:read,resource:write"
	DefaultTokenType = "user"
)

// ResetAccountChallenge creates a new confirmation token and resets the timer
func ResetAccountChallenge(ctx context.Context, account *Account) (*Account, error) {
	// FIXME add a mutex
	// FIXME this is crude!
	// FIXME deep-copy, no side effects!
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
	// FIXME deep-copy, no side effects!
	token, _ := util.ShortUUID()
	account.Expires = util.IncT(util.Timestamp(), DefaultAuthenticationExpiration)
	account.Ext2 = token
	account.Status = AccountLoggedOut

	if err := UpdateAccount(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func LogoutAccount(ctx context.Context, realm, clientID string) error {
	// FIXME add a mutex
	account, err := LookupAccount(ctx, realm, clientID)
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("account %s.%s not found", realm, clientID)
	}

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

func BlockAccount(ctx context.Context, realm, clientID string) error {
	// FIXME add a mutex
	account, err := LookupAccount(ctx, realm, clientID)
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("account %s.%s not found", realm, clientID)
	}

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

	account.Status = AccountBlocked
	return UpdateAccount(ctx, account)
}

// SendAccountChallenge sends a notification to the user promting to confirm the account
func SendAccountChallenge(ctx context.Context, account *Account) error {
	// FIXME use templates to send a proper email

	url := fmt.Sprintf("%s/login/%s", podops.DefaultAPIEndpoint, account.Ext1)

	if err := platform.SendEmail(env.GetString("EMAIL_FROM", "hello@podops.dev"), account.UserID, "Confirm your account", url); err != nil {
		return err
	}
	return nil
}

// SendAuthToken sends a notification to the user with the current authentication token
func SendAuthToken(ctx context.Context, account *Account) error {
	// FIXME this is not done, just a crude implementation

	if err := platform.SendEmail(env.GetString("EMAIL_FROM", "hello@podops.dev"), account.UserID, "Your confirmation token", account.Ext2); err != nil {
		return err
	}
	return nil

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
		return nil, http.StatusUnauthorized, nil
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

	account, err := FindAccountByUserID(ctx, req.Realm, req.UserID)
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
	auth, err = LookupAuthorization(ctx, account.Realm, account.ClientID)
	if err != nil {
		// FIXME maybe use a different code here
		return nil, http.StatusInternalServerError, err
	}
	if auth == nil {
		// FIXME this is hardcoded for podops, make it configurable
		auth = &Authorization{
			ClientID:  account.ClientID,
			Realm:     req.Realm,
			TokenType: DefaultTokenType,
			UserID:    req.UserID,
			Scope:     DefaultScope,
			Revoked:   false,
			Created:   now,
		}
	}

	auth.Token = createSimpleToken()
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

func createSimpleToken() string {
	token, _ := util.UUID()
	return token
}
