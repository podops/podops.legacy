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
		// content
		Title     string `json:"title"`
		Summary   string `json:"summary"`
		Published int64  `json:"published"`
		// metadata
		OriginURI    string `json:"origin"`    // e.g. the import URI
		Location     string `json:"location"`  // path to the backing resource file (.yaml,.mp3, etc.)
		EnclosureURI string `json:"enclosure"` // used in episode
		ImageURI     string `json:"image"`     // used in show, episode
		// internal
		Index   int   `json:"index"` // A running number that can be used to sort resources, e.g. episode number
		Created int64 `json:"-"`
		Updated int64 `json:"-"`
	}

	// ResourceList returns a list of resources
	ResourceList struct {
		Resources []*Resource `json:"resources" `
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
