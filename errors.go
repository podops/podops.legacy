package podops

import "errors"

var (
	// ErrNotAuthorized indicates that the API call is not authorized
	ErrNotAuthorized = errors.New("api: not authorized")
	// ErrNoToken indicates that no bearer token was provided
	ErrNoToken = errors.New("api: no token provided")

	// ErrInvalidParameters indicates that parameters used in an API call are not valid
	ErrInvalidParameters = errors.New("api: invalid parameters")
	// ErrValidationFailed indicates that a resource validation failed
	ErrValidationFailed = errors.New("api: validation failed")

	// ErrNoSuchProduction indicates that the production does not exist
	ErrNoSuchProduction = errors.New("api: production doesn't exist")
	// ErrNoSuchResource indicates that the resource does not exist
	ErrNoSuchResource = errors.New("api: resource doesn't exist")
	// ErrNoSuchAsset indicates that the asset does not exist
	ErrNoSuchAsset = errors.New("api: asset doesn't exist")
	// ErrBuildFailed indicates that the feed build failed
	ErrBuildFailed = errors.New("api: build failed")

	// ErrInternalError indicates that an unspecified internal error happened
	ErrInternalError = errors.New("api: internal error")
)
