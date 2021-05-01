package errordef

import "errors"

var (
	// ErrNotAuthorized indicates that the API caller is not authorized
	ErrNotAuthorized = errors.New("not authorized")
	// ErrNoToken indicates that no bearer token was provided
	ErrNoToken = errors.New("no token provided")
	// ErrNoToken indicates that the bearer token is not valid
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidRoute indicates that the route and/or its parameters are not valid
	ErrInvalidRoute = errors.New("invalid route")

	// ErrInvalidParameters indicates that parameters used in an API call are not valid
	ErrInvalidParameters = errors.New("invalid parameters")
	// ErrValidationFailed indicates that some validation failed
	ErrValidationFailed = errors.New("validation failed")

	// ErrNoSuch... indicates that the requested resource does not exist
	ErrNoSuchProduction = errors.New("production doesn't exist")
	ErrNoSuchEpisode    = errors.New("episode doesn't exist")
	ErrNoSuchAsset      = errors.New("asset doesn't exist")
	ErrNoSuchResource   = errors.New("resource doesn't exist")

	// ErrMissingResource indicates that a resource required for an operation can not be found
	ErrMissingResource = errors.New("can't find resource")

	// ErrBuildFailed indicates that there was an error while building the feed
	ErrBuildFailed = errors.New("build failed")
	// ErrFeedFailed indicates that some pre-requisites for building the feed are not met
	ErrFeedFailed = errors.New("can't build feed.xml")

	// ErrInvalidClientConfiguration indicates that the client configuration is in invalid
	ErrInvalidClientConfiguration = errors.New("invalid configuration")
	// ErrInvalidClientParameters indicates that parameters used in an client API call are not valid
	ErrInvalidClientParameters = errors.New("invalid parameters")

	// ErrInternalError indicates evrything else
	ErrInternalError = errors.New("internal error")
)
