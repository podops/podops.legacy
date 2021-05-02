package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/txsvc/platform/v2/pkg/id"
	"github.com/txsvc/platform/v2/pkg/timestamp"

	"github.com/podops/podops/internal/errordef"
	"github.com/podops/podops/internal/messagedef"
)

// ResetAccountChallenge creates a new confirmation token and resets the timer
func ResetAccountChallenge(ctx context.Context, account *Account) (*Account, error) {
	token, _ := id.ShortUUID()
	account.Expires = timestamp.IncT(timestamp.Now(), ac.authenticationExpiration)
	account.Ext1 = token
	account.Status = AccountUnconfirmed

	if err := UpdateAccount(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

// ResetAuthToken creates a new authorization token and resets the timer
func ResetAuthToken(ctx context.Context, account *Account) (*Account, error) {
	token, _ := id.ShortUUID()
	account.Expires = timestamp.IncT(timestamp.Now(), ac.authenticationExpiration)
	account.Ext2 = token
	account.Status = AccountLoggedOut

	if err := UpdateAccount(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func LogoutAccount(ctx context.Context, realm, clientID string) error {
	account, err := LookupAccount(ctx, realm, clientID)
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf(messagedef.MsgAuthenticationNotFound, fmt.Sprintf("%s.%s", realm, clientID))
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
	account, err := LookupAccount(ctx, realm, clientID)
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf(messagedef.MsgAuthenticationNotFound, fmt.Sprintf("%s.%s", realm, clientID))
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

// ConfirmLoginChallenge confirms the account
func ConfirmLoginChallenge(ctx context.Context, token string) (*Account, int, error) {
	if token == "" {
		return nil, http.StatusUnauthorized, errordef.ErrNoToken
	}

	account, err := FindAccountByToken(ctx, token)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if account == nil {
		return nil, http.StatusUnauthorized, nil
	}
	now := timestamp.Now()
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
	var auth *Authorization

	account, err := FindAccountByUserID(ctx, req.Realm, req.UserID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if account == nil {
		return nil, http.StatusNotFound, nil
	}
	now := timestamp.Now()
	if account.Expires < now || account.Ext2 != req.Token {
		return nil, http.StatusUnauthorized, nil
	}

	// all OK, create or update the authorization
	auth, err = LookupAuthorization(ctx, account.Realm, account.ClientID)
	if err != nil {
		return nil, http.StatusInternalServerError, err // FIXME maybe use a different code here
	}
	if auth == nil {
		auth = ac.createAuthorization(account, req)
	}
	auth.Token = CreateSimpleToken()
	auth.Expires = now + (int64(ac.authorizationExpiration) * 86400)
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
