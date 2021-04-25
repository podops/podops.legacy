package errordef

import "errors"

var (
	// ErrNotAuthorized indicates that the API caller is not authorized
	ErrNotAuthorized = errors.New("api: not authorized")
	// ErrNoToken indicates that no bearer token was provided
	ErrNoToken = errors.New("api: no token provided")
	// ErrNoToken indicates that the bearer token is not valid
	ErrInvalidToken = errors.New("api: invalid token")

	// ErrInvalidRoute indicates that the route and/or its parameters are not valid
	ErrInvalidRoute = errors.New("api: invalid route")

	// ErrInvalidParameters indicates that parameters used in an API call are not valid
	ErrInvalidParameters = errors.New("api: invalid parameters")
	// ErrValidationFailed indicates that some validation failed
	ErrValidationFailed = errors.New("api: validation failed")

	// ErrNoSuch... indicates that the requested resource does not exist
	ErrNoSuchProduction = errors.New("api: production doesn't exist")
	ErrNoSuchEpisode    = errors.New("api: episode doesn't exist")
	ErrNoSuchAsset      = errors.New("api: asset doesn't exist")
	ErrNoSuchResource   = errors.New("api: resource doesn't exist")

	// ErrMissingResource indicates that a resource required for an operation can not be found
	ErrMissingResource = errors.New("api: can't find resource")

	// ErrBuildFailed indicates that there was an error while building the feed
	ErrBuildFailed = errors.New("api: build failed")
	// ErrFeedFailed indicates that some pre-requisites for building the feed are not met
	ErrFeedFailed = errors.New("api: can't build feed.xml")

	// ErrInvalidClientConfiguration indicates that the client configuration is in invalid
	ErrInvalidClientConfiguration = errors.New("client: invalid configuration")
	// ErrInvalidClientParameters indicates that parameters used in an client API call are not valid
	ErrInvalidClientParameters = errors.New("client: invalid parameters")

	// ErrInternalError indicates evrything else
	ErrInternalError = errors.New("api: internal error")
)
