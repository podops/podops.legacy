package errordef

import "errors"

var (
	// ErrInvalidConfiguration indicates that the client configuration is in invalid
	ErrInvalidConfiguration = errors.New("client: invalid configuration")
	// ErrInvalidClientParameters indicates that parameters used in an client API call are not valid
	ErrInvalidClientParameters = errors.New("client: invalid parameters")
)
