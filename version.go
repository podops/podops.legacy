package podops

import "fmt"

/*
Make sure to also change here:
- platform/tasks.go
- podops-infra/inventory/*
*/

const (
	// Version specifies the verion of the API and its structs
	Version = "v1"

	// MajorVersion of the API
	MajorVersion = 0
	// MinorVersion of the API
	MinorVersion = 9
	// FixVersion of the API
	FixVersion = 13
)

var (
	// VersionString is the canonical API description
	VersionString string = fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, FixVersion)
	// UserAgentString identifies any http request podops makes
	UserAgentString string = fmt.Sprintf("PodOps %d.%d.%d", MajorVersion, MinorVersion, FixVersion)
)
