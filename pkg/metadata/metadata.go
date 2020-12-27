package metadata

import (
	"time"
)

const (
	// DefaultPortalEndpoint points to the portal
	DefaultPortalEndpoint = "https://podops.dev"
	// DefaultAPIEndpoint points to the API
	DefaultAPIEndpoint = "https://api.podops.dev"
	// DefaultCDNEndpoint point to the CDN
	DefaultCDNEndpoint = "https://cdn.podops.dev"

	//
	// Required and optional labels:
	//
	//	show:
	//		language:	<ISO639 two-letter-code> REQUIRED 'channel.language'
	//		explicit:	True | False REQUIRED 'channel.itunes.explicit'
	//		type:		Episodic | Serial REQUIRED 'channel. itunes.type'
	//		block:		Yes OPTIONAL 'channel.itunes.block' Anything else than 'Yes' has no effect
	//		complete:	Yes OPTIONAL 'channel.itunes.complete' Anything else than 'Yes' has no effect
	//
	//	episode:
	//		guid:		<unique id> 'item.guid'
	//		date:		<publish date> REQUIRED 'item.pubDate'
	//		season: 	<season number> OPTIONAL 'item.itunes.season'
	//		episode:	<episode number> REQUIRED 'item.itunes.episode'
	//		explicit:	True | False REQUIRED 'channel.itunes.explicit'
	//		type:		Full | Trailer | Bonus REQUIRED 'item.itunes.episodeType'
	//		block:		Yes OPTIONAL 'item.itunes.block' Anything else than 'Yes' has no effect
	//

	LabelLanguage   = "language"
	LabelExplicit   = "explicit"
	LabelType       = "type"
	LabelBlock      = "block"
	LabelComplete   = "complete"
	LabelGUID       = "guid"
	LabelParentGUID = "parent_guid"
	LabelDate       = "date"
	LabelSeason     = "season"
	LabelEpisode    = "episode"

	ShowTypeEpisodic = "Episodic"
	ShowTypeSerial   = "Serial"

	EpisodeTypeFull    = "Full"
	EpisodeTypeTrailer = "Trailer"
	EpisodeTypeBonus   = "Bonus"
)

type (
	// Apple Podcast: https://help.apple.com/itc/podcasts_connect/#/itcb54353390
	// RSS 2.0: https://cyber.harvard.edu/rss/rss.html

	// Metadata contains information describing a resource
	Metadata struct {
		Name   string            `json:"name" yaml:"name" binding:"required"` // REQUIRED <unique name>
		Labels map[string]string `json:"labels" yaml:"labels,omitempty"`      // REQUIRED
	}

	// ResourceMetadata holds only the kind and metadata a resource
	ResourceMetadata struct {
		APIVersion string   `json:"apiVersion" yaml:"apiVersion" binding:"required"` // REQUIRED default: v1.0
		Kind       string   `json:"kind" yaml:"kind" binding:"required"`             // REQUIRED default: show
		Metadata   Metadata `json:"metadata" yaml:"metadata" binding:"required"`     // REQUIRED
	}

	// Show holds all metadata related to a podcast/show
	Show struct {
		APIVersion  string          `json:"apiVersion" yaml:"apiVersion" binding:"required"`   // REQUIRED default: v1.0
		Kind        string          `json:"kind" yaml:"kind" binding:"required"`               // REQUIRED default: show
		Metadata    Metadata        `json:"metadata" yaml:"metadata" binding:"required"`       // REQUIRED
		Description ShowDescription `json:"description" yaml:"description" binding:"required"` // REQUIRED
		Image       Resource        `json:"image" yaml:"image" binding:"required"`             // REQUIRED 'channel.itunes.image'
	}

	// Episode holds all metadata related to a podcast episode
	Episode struct {
		APIVersion  string             `json:"apiVersion" yaml:"apiVersion" binding:"required"`   // REQUIRED default: v1.0
		Kind        string             `json:"kind" yaml:"kind" binding:"required"`               // REQUIRED default: episode
		Metadata    Metadata           `json:"metadata" yaml:"metadata" binding:"required"`       // REQUIRED
		Description EpisodeDescription `json:"description" yaml:"description" binding:"required"` // REQUIRED
		Image       Resource           `json:"image" yaml:"image" binding:"required"`             // REQUIRED 'item.itunes.image'
		Enclosure   Resource           `json:"enclosure" yaml:"enclosure" binding:"required"`     // REQUIRED
	}

	// ShowDescription holds essential show metadata
	ShowDescription struct {
		Title     string    `json:"title" yaml:"title" binding:"required"`          // REQUIRED 'channel.title' 'channel.itunes.title'
		Summary   string    `json:"summary" yaml:"summary" binding:"required"`      // REQUIRED 'channel.description'
		Link      Resource  `json:"link" yaml:"link"`                               // RECOMMENDED 'channel.link'
		Category  Category  `json:"category" yaml:"category" binding:"required"`    // REQUIRED channel.category
		Owner     Owner     `json:"owner" yaml:"owner"`                             // RECOMMENDED 'channel.itunes.owner'
		Author    string    `json:"author" yaml:"author"`                           // RECOMMENDED 'channel.itunes.author'
		Copyright string    `json:"copyright,omitempty" yaml:"copyright,omitempty"` // OPTIONAL 'channel.copyright'
		NewFeed   *Resource `json:"newFeed,omitempty" yaml:"newFeed,omitempty"`     // OPTIONAL channel.itunes.new-feed-url -> move to label
	}

	// EpisodeDescription holds essential episode metadata
	EpisodeDescription struct {
		Title       string   `json:"title" yaml:"title" binding:"required"`                                 // REQUIRED 'item.title' 'item.itunes.title'
		Summary     string   `json:"summary" yaml:"summary" binding:"required"`                             // REQUIRED 'item.description'
		EpisodeText string   `json:"episodeText,omitempty" yaml:"episodeText,omitempty" binding:"required"` // REQUIRED 'item.itunes.summary'
		Link        Resource `json:"link" yaml:"link"`                                                      // RECOMMENDED 'item.link'
		Duration    int      `json:"duration" yaml:"duration" binding:"required"`                           // REQUIRED 'item.itunes.duration'
	}

	// Owner describes the owner of the show/podcast
	Owner struct {
		Name  string `json:"name" yaml:"name" binding:"required"`   // REQUIRED
		Email string `json:"email" yaml:"email" binding:"required"` // REQUIRED
	}

	// Category is the show/episodes category and it's subcategories
	Category struct {
		Name        string   `json:"name" yaml:"name" binding:"required"`      // REQUIRED
		SubCategory []string `json:"subcategory" yaml:"subcategory,omitempty"` // OPTIONAL
	}

	// Resource provides a link to a media resource
	Resource struct {
		URI    string `json:"uri" yaml:"uri" binding:"required"`        // REQUIRED
		Title  string `json:"title,omitempty" yaml:"title,omitempty"`   // OPTIONAL
		Anchor string `json:"anchor,omitempty" yaml:"anchor,omitempty"` // OPTIONAL
		Rel    string `json:"rel,omitempty" yaml:"rel,omitempty"`       // OPTIONAL
		Type   string `json:"type,omitempty" yaml:"type,omitempty"`     // OPTIONAL
		Size   int    `json:"size,omitempty" yaml:"size,omitempty"`     // OPTIONAL
	}
)

//
// Some helper functions to deal with metadata
//

// PublishDateTimestamp converts a RFC1123Z formatted timestamp into UNIX timestamp
func (e *Episode) PublishDateTimestamp() int64 {
	pd := e.Metadata.Labels[LabelDate]
	if pd == "" {
		return 0
	}
	t, err := time.Parse(time.RFC1123Z, pd)
	if err != nil {
		return 0
	}

	return t.Unix()
}

// PublishDate is a convenience method to access the pub date
func (e *Episode) PublishDate() string {
	return e.Metadata.Labels[LabelDate]
}

// GUID is a convenience method to access the resources guid
func (e *Episode) GUID() string {
	return e.Metadata.Labels[LabelGUID]
}

// GUID is a convenience method to access the resources guid
func (r *ResourceMetadata) GUID() string {
	return r.Metadata.Labels[LabelGUID]
}

// GUID is a convenience method to access the resources guid
func (s *Show) GUID() string {
	return s.Metadata.Labels[LabelGUID]
}
