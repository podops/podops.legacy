package auth

type (

	// OldAuthorizationRequest struct is used to request a token
	OldAuthorizationRequest struct {
		Secret     string `json:"secret" binding:"required"`
		Realm      string `json:"realm" binding:"required"`
		ClientID   string `json:"client_id" binding:"required"`
		ClientType string `json:"client_type" binding:"required"` // user,app,bot
		UserID     string `json:"user_id" binding:"required"`
		Scope      string `json:"scope" binding:"required"`
		Duration   int64  `json:"duration" binding:"required"`
	}
)
