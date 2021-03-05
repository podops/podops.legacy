package auth

import (
	"github.com/fupas/commons/pkg/util"
)

type (

	// Client represents the claim of the client calling the API
	Client struct {
		ClientID string `json:"client_id"`
		UserID   string `json:"user_id"`
		Scope    string `json:"scope"`
	}
)

// CreateJWTToken creates a token that can be used for JWT authentication / authorization
func CreateJWTToken(secret, realm, clientID, userID, scope string, duration int64) (string, error) {

	token, _ := util.UUID()
	return token, nil
}
