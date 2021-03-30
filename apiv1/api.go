package apiv1

const (
	// NamespacePrefix namespace for the client and CLI
	NamespacePrefix = "/a/v1"
	// GraphqlNamespacePrefix namespace for the GraphQL endpoints
	GraphqlNamespacePrefix = "/q"
	// AdminNamespacePrefix namespace for internal admin endpoints
	AdminNamespacePrefix = "/_a"
	// TaskNamespacePrefix namespace for internal Cloud Task callbacks
	TaskNamespacePrefix = "/_t"
	// ContentNamespace namespace for thr CDN
	ContentNamespace = "/c"

	// All the API & CLI endpoint routes

	// LoginRequestRoute route to LoginRequestEndpoint
	LoginRequestRoute = "/login"
	// LogoutRequestRoute route to LogoutRequestEndpoint
	LogoutRequestRoute = "/logout"
	// LoginConfirmationRoute route to LoginConfirmationEndpoint
	LoginConfirmationRoute = "/login/:token"
	// GetAuthorizationRoute route to GetAuthorizationEndpoint
	GetAuthorizationRoute = "/auth"

	// ProductionRoute route to ProductionEndpoint
	ProductionRoute = "/production"
	// ListProductionsRoute route to ListProductionsEndpoint
	ListProductionsRoute = "/productions"

	// FindResourceRoute route to FindResourceEndpoint
	FindResourceRoute = "/resource/:id"
	// GetResourceRoute route to ResourceEndpoint
	GetResourceRoute = "/resource/:prod/:kind/:id"
	// ListResourcesRoute route to ResourceEndpoint GET
	ListResourcesRoute = "/resource/:prod/:kind"
	// UpdateResourceRoute route to ResourceEndpoint POST,PUT
	UpdateResourceRoute = "/resource/:prod/:kind/:id"
	// DeleteResourceRoute route to ResourceEndpoint
	DeleteResourceRoute = "/resource/:prod/:kind/:id"

	// ImportTask route to ImportTaskEndpoint
	ImportTask = "/import"
	// full canonical route
	ImportTaskWithPrefix = "/_t/import"

	// BuildRoute route to BuildEndpoint
	BuildRoute = "/build"
	// UploadRoute route to UploadEndpoint
	UploadRoute = "/upload/:prod"

	// ShowRoute route to show.json
	ShowRoute = "/s/:name"

	// EpisodeRoute route to show.json
	EpisodeRoute = "/e/:guid"

	// FeedRoute route to feed.xml
	FeedRoute = "/s/:name/feed.xml"

	// DefaultCDNRoute route to /:guid/:asset
	DefaultCDNRoute = "/:guid/:asset"

	// GraphqlRoute route to GraphqlEndpoint
	GraphqlRoute = "/query"

	// GraphqlPlaygroundRoute route to GraphqlPlaygroundEndpoint
	GraphqlPlaygroundRoute = "/playground"
)
