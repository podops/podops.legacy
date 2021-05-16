package messagedef

const (
	MsgStatus      = "status: %d"
	MsgClientError = "internal client error"

	MsgResourceNotFound              = "resource '%s' not found"
	MsgResourceAlreadyExists         = "resource '%s' already exists"
	MsgResourceIsInvalid             = "invalid resource '%s'"
	MsgResourceInconsistentInventory = "resource count mismatch: expected %d, got %d. '%s,%s'"
	MsgResourceKindMismatch          = "resource mismatch: expected '%s', got '%s'"
	MsgResourceUnsupportedKind       = "unsupported kind '%s'"
	MsgResourceInvalidGUID           = "invalid guid '%s'"

	MsgResourceImportError = "error transfering '%s'"
	MsgResourceUploadError = "error uploading '%s'"

	MsgParameterIsInvalid = "invalid parameter '%s'"
	MsgParameterMismatch  = "parameters mismatch. expected '%s', got '%s'"

	MsgAuthenticationNotFound     = "account '%s' not found"
	MsgAuthenticationTokenExpired = "token expired"
	MsgAuthenticationTokenInvalid = "token is invalid"

	// CLI

	MsgLoginNewAccount   = "New account created. Check your inbox and confirm the email address."
	MsgLoginVerification = "Login verificaction sent. Check your inbox."
	MsgLoginInvalidEmail = "'%s' is not a valid email address"
	MsgLoginError        = "already logged-in"
	MsgLoginSuccess      = "Login successful"
	MsgLogoutSuccess     = "Logout successful."
	MsgNotLoggedIn       = "Not logged in."
	MsgServerError       = "something went wrong: [%d]"

	MsgArgumentMissing       = "missing argument '%s'"
	MsgTooManyArguments      = "too many arguments"
	MsgArgumentCountMismatch = "argument mismatch: expected %d, got %d"

	MsgResourceCreated       = "created resource '%s'"
	MsgResourceUpdated       = "updated resource '%s'"
	MsgResourceDeleted       = "deleted resource '%s'"
	MsgResourceUnknown       = "unknown resource '%s'"
	MsgResourceDeletingError = "error deleting resource '%s'"
	MsgResourceUploadSuccess = "uploaded '%s'"

	MsgNoProductionsFound = "production(s) not found"
	MsgNoResourcesFound   = "resource(s) not found"

	MsgErrorUpdatingConfig      = "error updating config."
	MsgErrorNoProduction        = "no production set. Use 'po show [ID|name]' first"
	MsgErrorCanNotSetProduction = "no production set. Use 'po shows' to find available productions"

	MsgBuildSuccess = "build production '%s' successful.\nAccess the feed at %s"
)
