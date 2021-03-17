package auth

import (
	"context"
	"fmt"
	"net/http"
)

const (
	// DatastoreAuthorizations collection AUTHORIZATION
	DatastoreAuthorizations string = "AUTHORIZATIONS"
	// DatastoreAccounts collection ACCOUNTS
	DatastoreAccounts string = "ACCOUNTS"

	// AuthTypeSimpleToken constant token
	AuthTypeSimpleToken = "token"
	// AuthTypeJWT constant jwt
	AuthTypeJWT = "jwt"
	// AuthTypeSlack constant slack
	AuthTypeSlack = "slack"
)

// LookupAccount retrieves an account within a given realm
func LookupAccount(ctx context.Context, realm, userID string) (*Account, error) {
	return nil, fmt.Errorf("LookupAccount: not implemented")
}

// CreateAccount creates an new account within a given realm
func CreateAccount(ctx context.Context, realm, userID string) (*Account, error) {
	return nil, fmt.Errorf("CreateAccount: not implemented")
}

// ResetAccountChallenge creates a new confirmation token and resets the timer
func ResetAccountChallenge(ctx context.Context, account *Account) (*Account, error) {
	return nil, fmt.Errorf("ResetLoginChallenge: not implemented")
}

// SendAccountChallenge sends a notification to the user promting to confirm the account
func SendAccountChallenge(ctx context.Context, account *Account) error {
	return fmt.Errorf("SendLoginChallenge: not implemented")
}

// ConfirmLoginChallenge confirms the account
func ConfirmLoginChallenge(ctx context.Context, token string) (*Account, int, error) {
	return nil, http.StatusInternalServerError, fmt.Errorf("ConfirmLoginChallenge: not implemented")
}

// ResetAuthToken creates a new authorization token and resets the timer
func ResetAuthToken(ctx context.Context, account *Account) (*Account, error) {
	return nil, fmt.Errorf("ResetAuthToken: not implemented")
}

// SendAuthToken sends a notification to the user with the current authentication token
func SendAuthToken(ctx context.Context, account *Account) error {
	return fmt.Errorf("SendAuthToken: not implemented")
}

// CreateAuthentication confirms the temporary auth token and creates the permanent one
func CreateAuthentication(ctx context.Context, req *AuthorizationRequest) (*Account, int, error) {
	return nil, http.StatusInternalServerError, fmt.Errorf("CreateAuthentication: not implemented")
}
