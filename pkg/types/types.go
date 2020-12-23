package types

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
)
