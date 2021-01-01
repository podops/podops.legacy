package apiv1

import (
	"fmt"
)

type (
	// ProductionRequest defines the request
	ProductionRequest struct {
		Name    string `json:"name" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Summary string `json:"summary" binding:"required"`
		GUID    string `json:"guid,omitempty"`
	}

	// ProductionResponse defines the request
	ProductionResponse struct {
		Name string `json:"name" binding:"required"`
		GUID string `json:"guid" binding:"required"`
	}

	// ProductionsResponse returns a list of productions
	ProductionsResponse struct {
		List []ProductionDetails `json:"list" `
	}

	// ProductionDetails provides details about a production
	ProductionDetails struct {
		Name  string `json:"name" binding:"required"`
		GUID  string `json:"guid" binding:"required"`
		Title string `json:"title,omitempty" `
	}

	// BuildRequest initiates the build of the feed
	BuildRequest struct {
		GUID string `json:"guid" binding:"required"`
	}

	// BuildResponse returns the resulting URL
	BuildResponse struct {
		GUID string `json:"guid" binding:"required"`
		URL  string `json:"url" binding:"required"`
	}

	// StatusObject is used to report status and errors in an API request.
	// The struct can be used as a response object or be treated as an error object
	StatusObject struct {
		Status  int    `json:"status" binding:"required"`
		Message string `json:"message" binding:"required"`
	}

	// ImportRequest is used by the import task
	ImportRequest struct {
		Source string `json:"src" binding:"required"`
		Dest   string `json:"dest" binding:"required"`
	}

	// AuthorizationRequest struct is used to request a token
	// Imported from https://github.com/txsvc/service/blob/main/pkg/auth/types.go
	AuthorizationRequest struct {
		Secret     string `json:"secret" binding:"required"`
		Realm      string `json:"realm" binding:"required"`
		ClientID   string `json:"client_id" binding:"required"`
		ClientType string `json:"client_type" binding:"required"` // user,app,bot
		UserID     string `json:"user_id" binding:"required"`
		Scope      string `json:"scope" binding:"required"`
		Duration   int64  `json:"duration" binding:"required"`
	}

	// AuthorizationResponse provides a valid token
	// Imported from https://github.com/txsvc/service/blob/main/pkg/auth/types.go
	AuthorizationResponse struct {
		Realm    string `json:"realm" binding:"required"`
		ClientID string `json:"client_id" binding:"required"`
		Token    string `json:"token" binding:"required"`
	}
)

func (so *StatusObject) Error() string {
	return fmt.Sprintf("%s: %d", so.Message, so.Status)
}
