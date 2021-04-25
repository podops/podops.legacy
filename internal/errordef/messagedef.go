package errordef

const (
	MsgStatus      = "status: %d"
	MsgCLIError    = "internal cli error"
	MsgClientError = "internal client error"

	MsgResourceNotFound      = "resource '%s' not found"
	MsgResourceAlreadyExists = "resource '%s' already exists"
	MsgInvalidResource       = "invalid resource '%s'"
	MsgResourceCountMismatch = "resource count mismatch: expected %d, got %d"
	MsgResourceMismatch      = "parameters mismatch: expected '%s', got '%s'"

	MsgErrUploadingResource = "error transfering '%s'"
	MsgErrSyncingResource   = "error syncing '%s'"

	MsgArgumentCountMismatch = "args: expected %d, got %d"

	MsgInvalidGUID        = "invalid guid '%s'"
	MsgInvalidParameter   = "invalid parameter '%s'"
	MsgInvalidIdentifier  = "invalid identifier '%s'"
	MsgParametersMismatch = "parameters mismatch. expected '%s', got '%s'"
	MsgUnsupportedType    = "unsupported type '%s'"

	MsgAccountNotFound = "account '%s' not found"
)
