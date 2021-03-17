package auth

const (
	// AccountActive indicates a confirmed account with a valid login
	AccountActive = 1
	// AccountLoggedOut indicates a confirmed account without a valid login
	AccountLoggedOut = 0
	// AccountDeactivated indicates an account that has been deactivated due to
	// e.g. account deletion or UserID swap
	AccountDeactivated = -1
	// AccountBlocked signals an issue with the account that needs intervention
	AccountBlocked = -2
	// AccountUnconfirmed well guess what?
	AccountUnconfirmed = -3
)

type (
	// AuthorizationRequest represents a login/authorization request from a user, app, or bot
	AuthorizationRequest struct {
		Realm    string `json:"realm" binding:"required"`
		UserID   string `json:"user_id" binding:"required"`
		ClientID string `json:"client_id"`
		Token    string `json:"token"`
	}

	// Account represents an account for a user or client (e.g. API, bot)
	Account struct {
		Realm    string `json:"realm"`
		ClientID string `json:"client_id"` // a unique id within [realm,user_id]
		UserID   string `json:"user_id"`   // external id for the entity e.g. email for a user
		// status and other metadata
		Status int `json:"status"` // default == AccountUnconfirmed
		// internal
		TempToken string `json:"-"` // used to confirm the account and then to request the real token
		Expires   int64  `json:"-"` // 0 == never
		Confirmed int64  `json:"-"`
		Created   int64  `json:"-"`
		Updated   int64  `json:"-"`
	}
)
