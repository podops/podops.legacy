package types

import "fmt"

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

	// StatusObject is used to report status and errors in an API request.
	// The struct can be used as a response object or be treated as an error object
	StatusObject struct {
		Status  int    `json:"status" binding:"required"`
		Message string `json:"message" binding:"required"`
	}
)

func (so *StatusObject) Error() string {
	return fmt.Sprintf("%s: %d", so.Message, so.Status)
}