package podops

import (
	"fmt"
)

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
		// Metadata used in e.g. the web UI
		Title     string `json:"title"`
		Summary   string `json:"summary"`
		Extra1    string `json:"extra1"` // These two attributes are just placeholders for any kind of resource specific data
		Extra2    string `json:"extra2"` // One possible use is to e.g. store the URL of an episodes media file here
		Published int64  `json:"published"`
		Location  string `json:"location"` // path to the .yaml
		// Media metadata used for e.g. .mp3/.png
		Image       string `json:"image"`        // Full URL to the show/episode image
		ContentType string `json:"content_type"` // DEPRECATED move to ResourceMetadata
		Duration    int64  `json:"duration"`     // DEPRECATED move to ResourceMetadata
		Size        int64  `json:"size"`         // DEPRECATED move to ResourceMetadata
		// internal
		Index   int   `json:"index"` // A running number that can be used to sort resources, e.g. episode number
		Created int64 `json:"-"`
		Updated int64 `json:"-"`
	}

	// ResourceList returns a list of resources
	ResourceList struct {
		Resources []*Resource `json:"resources" `
	}

	// ContentMetadata keeps basic metadata of a cdn resource
	ContentMetadata struct {
		Name        string `json:"name"`
		GUID        string `json:"guid"`
		Size        int64  `json:"size"`
		Duration    int64  `json:"duration"`
		ContentType string `json:"content_type"`
		Etag        string `json:"etag"`
		Timestamp   int64  `json:"timestamp"`
	}

	// BuildRequest initiates the build of the feed
	BuildRequest struct {
		GUID         string `json:"guid" binding:"required"`
		FeedURL      string `json:"feed"`
		FeedAliasURL string `json:"alias"`
	}

	// ImportRequest is used by the import task
	ImportRequest struct {
		GUID     string `json:"guid" binding:"required"`
		Source   string `json:"src" binding:"required"`
		Original string `json:"original" binding:"required"`
		Dest     string `json:"dest"`
	}
)

//
// helper functions to work with the models
//

// GetPublicLocation returns the public url of a resource if it exists on the CDN or an empty string otherwise
func (r *Resource) GetPublicLocation() string {
	if r.Kind == ResourceAsset {
		return fmt.Sprintf("%s/%s", DefaultStorageEndpoint, r.Location)
	}
	return ""
}
