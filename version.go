package podops

import "fmt"

/*
Make sure to update version numbers in these locations also:

- internal/transport/transport.go
- .github/*
*/

const (
	// Version specifies the verion of the API and its structs
	Version = "v1"

	// MajorVersion of the API
	MajorVersion = 1
	// MinorVersion of the API
	MinorVersion = 0
	// FixVersion of the API
	FixVersion = 2
)

var (
	// VersionString is the canonical API description
	VersionString string = fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, FixVersion)
)
