package auth

import "github.com/fupas/commons/pkg/util"

type (
	// Authorization represents a user, app or bot and its permissions
	Authorization struct {
		ClientID  string `json:"client_id" binding:"required"` // UNIQUE
		Name      string `json:"name"`                         // name of the domain, realm, tennant etc
		Token     string `json:"token" binding:"required"`
		TokenType string `json:"token_type" binding:"required"` // user,app,bot
		UserID    string `json:"user_id"`                       // depends on TokenType. UserID could equal ClientID or BotUSerID in Slack
		Scope     string `json:"scope"`                         // a comma separated list of scopes, see below
		Expires   int64  `json:"expires"`                       // 0 = never
		// internal
		// FIXME: add revokation flag to the Authorization
		AuthType string `json:"-"` // currently: jwt, slack
		Created  int64  `json:"-"`
		Updated  int64  `json:"-"`
	}

	// AuthorizationRequest struct is used to request a token
	AuthorizationRequest struct {
		Secret     string `json:"secret" binding:"required"`
		Realm      string `json:"realm" binding:"required"`
		ClientID   string `json:"client_id" binding:"required"`
		ClientType string `json:"client_type" binding:"required"` // user,app,bot
		UserID     string `json:"user_id" binding:"required"`
		Scope      string `json:"scope" binding:"required"`
		Duration   int64  `json:"duration" binding:"required"`
	}

	// AuthorizationResponse provides the token to the requestor
	AuthorizationResponse struct {
		Realm    string `json:"realm" binding:"required"`
		ClientID string `json:"client_id" binding:"required"`
		Token    string `json:"token" binding:"required"`
	}
)

// IsValid verifies that the Authorization is still valid, i.e. not expired and not revoked.
func (a *Authorization) IsValid() bool {
	if a.Expires == 0 {
		return true
	}
	if a.Expires > util.Timestamp() {
		return true
	}
	return false
}
