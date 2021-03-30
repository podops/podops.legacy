package podops

import (
	"fmt"
)

type (

	// FIXME find a better place for StatusObject

	// StatusObject is used to report operation status and errors in an API request.
	// The struct can be used as a response object or be treated as an error object
	StatusObject struct {
		Status    int    `json:"status" binding:"required"`
		Message   string `json:"message" binding:"required"`
		RootError error  `json:"-"`
	}

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

	// BuildRequest initiates the build of the feed
	BuildRequest struct {
		GUID         string `json:"guid" binding:"required"`
		FeedURL      string `json:"feed"`
		FeedAliasURL string `json:"alias"`
	}

	// ImportRequest is used by the import task
	ImportRequest struct {
		Source   string `json:"src" binding:"required"`
		Dest     string `json:"dest" binding:"required"`
		Original string `json:"original" binding:"required"`
	}
)

//
// helper functions to work with the models
//

// NewStatus initializes a new StatusObject
func NewStatus(s int, m string) StatusObject {
	return StatusObject{Status: s, Message: m}
}

// NewErrorStatus initializes a new StatusObject from an error
func NewErrorStatus(s int, e error) StatusObject {
	return StatusObject{Status: s, Message: e.Error(), RootError: e}
}

func (so *StatusObject) Error() string {
	return fmt.Sprintf("%s: %d", so.Message, so.Status)
}

// GetPublicLocation returns the public url of a resource if it exists on the CDN or an empty string otherwise
func (r *Resource) GetPublicLocation() string {
	if r.Kind == ResourceAsset {
		return fmt.Sprintf("%s/c/%s", DefaultCDNEndpoint, r.Location)
	}
	return ""
}