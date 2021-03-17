package apiv1

type (
	// Production is the parent struct of all other resources.
	Production struct {
		Name      string `json:"name" binding:"required"`
		GUID      string `json:"guid,omitempty"`
		Owner     string `json:"owner"`
		Title     string `json:"title"`
		Summary   string `json:"summary"`
		BuildDate int64  `json:"build_date"`
		// internal
		Created int64 `json:"-"`
		Updated int64 `json:"-"`
	}

	// ProductionList returns a list of productions
	ProductionList struct {
		Productions []*Production `json:"productions" `
	}

	// Resource is used to maintain a repository of all existing resources across all shows
	Resource struct {
		Name       string `json:"name"`
		GUID       string `json:"guid"`
		Kind       string `json:"kind"`
		ParentGUID string `json:"parent_guid"`
		Location   string `json:"location"` // path to the .yaml
		// Metadata used in e.g. the web UI
		Title     string `json:"title"`
		Summary   string `json:"summary"`
		Published int64  `json:"published"`
		Index     int    `json:"index"`  // A running number that can be used to sort resources, e.g. episode number
		Extra1    string `json:"extra1"` // These two attributes are just placeholders for any kind of resource type specific data.
		Extra2    string `json:"extra2"` // One possible use is to e.g. store the URL of an episodes media file here.
		// Media metadata used for e.g. .mp3/.png
		Image       string `json:"image"` // Full URL to the show/episode image
		ContentType string `json:"content_type"`
		Duration    int64  `json:"duration"`
		Size        int64  `json:"size"`
		// internal
		Created int64 `json:"-"`
		Updated int64 `json:"-"`
	}

	// ResourceList returns a list of resources
	ResourceList struct {
		Resources []*Resource `json:"resources" `
	}

	// Build initiates the build of the feed
	Build struct {
		GUID         string `json:"guid" binding:"required"`
		FeedURL      string `json:"feed"`
		FeedAliasURL string `json:"alias"`
	}

	// Import is used by the import task
	Import struct {
		Source string `json:"src" binding:"required"`
		Dest   string `json:"dest" binding:"required"`
	}

	// OldAuthorizationRequest struct is used to request a token
	// Imported from https://github.com/txsvc/service/blob/main/pkg/auth/types.go
	OldAuthorizationRequest struct {
		Secret     string `json:"secret" binding:"required"`
		Realm      string `json:"realm" binding:"required"`
		ClientID   string `json:"client_id" binding:"required"`
		ClientType string `json:"client_type" binding:"required"` // user,app,bot
		UserID     string `json:"user_id" binding:"required"`
		Scope      string `json:"scope" binding:"required"`
		Duration   int64  `json:"duration" binding:"required"`
	}

	// OldAuthorizationResponse provides a valid token
	// Imported from https://github.com/txsvc/service/blob/main/pkg/auth/types.go
	OldAuthorizationResponse struct {
		Realm    string `json:"realm" binding:"required"`
		ClientID string `json:"client_id" binding:"required"`
		Token    string `json:"token" binding:"required"`
	}
)
