package api

const (
	// AdminNamespacePrefix namespace for internal admin endpoints
	AdminNamespacePrefix = "/_a"
	// NamespacePrefix namespace for the CLI. Should not be used directly.
	NamespacePrefix = "/a/v1"

	// All the API & CLI endpoint routes

	// ProductionRoute route to ProductionEndpoint
	ProductionRoute = "/new"
	// ResourceRoute route to ResourceEndpoint
	ResourceRoute = "/update/:rsrc/:id"
)
