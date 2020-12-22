package cli

const (
	CmdLineName    = "po"
	CmdLineVersion = "v0.1"

	BasicCmdGroup    = "Basic Commands"
	SettingsCmdGroup = "Settings Commands"
	ShowCmdGroup     = "Show Commands"
	ShowMgmtCmdGroup = "Show Management Commands"

	// All the API & CLI endpoint routes

	// NewShowRoute creates a new production
	NewShowRoute = "/new"
	// CreateRoute creates a resource
	CreateRoute = "/create/:id/:rsrc"
	// UpdateRoute updates a resource
	UpdateRoute = "/update/:id/:rsrc"
)
