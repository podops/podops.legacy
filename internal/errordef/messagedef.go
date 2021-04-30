package errordef

const (
	MsgStatus      = "status: %d"
	MsgClientError = "internal client error"

	MsgResourceNotFound      = "resource '%s' not found"
	MsgResourceAlreadyExists = "resource '%s' already exists"
	MsgInvalidResource       = "invalid resource '%s'"
	MsgResourceCountMismatch = "resource count mismatch: expected %d, got %d"
	MsgResourceMismatch      = "parameters mismatch: expected '%s', got '%s'"

	MsgErrUploadingResource = "error transfering '%s'"
	MsgErrSyncingResource   = "error syncing '%s'"

	MsgArgumentCountMismatch = "args: expected %d, got %d"
	MsgMissingArgument       = "missing argument '%s'"

	MsgInvalidGUID        = "invalid guid '%s'"
	MsgInvalidParameter   = "invalid parameter '%s'"
	MsgInvalidIdentifier  = "invalid identifier '%s'"
	MsgParametersMismatch = "parameters mismatch. expected '%s', got '%s'"
	MsgUnsupportedType    = "unsupported type '%s'"

	MsgAccountNotFound = "account '%s' not found"

	// command line interface messages
	MsgCLIError                 = "internal cli error"
	MsgCLIStatus                = "status: %d"
	MsgCLIArgumentCountMismatch = "args: expected %d, got %d"
	MsgCLIProductionsNotFound   = "production(s) not found"
	MsgCLIResourcesNotFound     = "resource(s) not found"
	MsgCLIResourceUnknown       = "unknown resource type '%s'"
	MsgCLIResourceCreated       = "created resource '%s'"
	MsgCLIResourceUpdated       = "updated resource '%s'"
	MsgCLIResourceDeleted       = "deleted resource '%s'"
	MsgCLIErrorDeletingResource = "error deleting resource '%s'"
	MsgCLIUploadedResource      = "uploaded '%s'"

	MsgCLINewAccount        = "New account created. Check your inbox and confirm the email address."
	MsgCLILoginVerification = "Login verificaction sent. Check your inbox."
	MsgCLILoginError        = "Already logged-in, use 'po logout' first."
	MsgCLIErrorUpdateConfig = "Error updating config."
	MsgCLIAuthSuccess       = "Sucessfully authenticated."
	MsgCLITokenExpired      = "Token is expired"
	MsgCLITokenInvalid      = "Invalid token"
	MsgCLILogout            = "Logout successful."

	MsgCLIErrorNoProductionSet = "No production set. Use 'po show [ID|name]' first"
	MsgCLIErrorCanNotSet       = "No production set. Use 'po shows' to find available productions"
	MsgCLIBuild                = "Build production '%s' successful.\nAccess the feed at %s"
)