package podops

import (
	"fmt"
	"strings"
	"time"

	"github.com/txsvc/platform/v2/pkg/id"

	"github.com/podops/podops/internal/metadata"
)

const (
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

	// LabelLanguage ISO-639 two-letter language code. channel.language
	LabelLanguage = "language"
	// LabelExplicit ["true"|"false"] channel.itunes.explicit
	LabelExplicit = "explicit"
	// LabelType ["Episodic"|"Serial"] channel.itunes.type
	LabelType = "type"
	// LabelBlock ["Yes"] channel.itunes.block
	LabelBlock = "block"
	// LabelComplete ["Yes"] channel.itunes.complete
	LabelComplete = "complete"
	// LabelGUID resources GUID
	LabelGUID = "guid"
	// LabelParentGUID guid of the resources parent resource
	LabelParentGUID = "parent_guid"
	// LabelDate used as e.g. publish date of an episode
	LabelDate = "date"
	// LabelSeason defaults to "1"
	LabelSeason = "season"
	// LabelEpisode positive integer 1..
	LabelEpisode = ResourceEpisode

	// ShowTypeEpisodic type of podcast is episodic
	ShowTypeEpisodic = "Episodic"
	// ShowTypeSerial type of podcast is serial
	ShowTypeSerial = "Serial"

	// EpisodeTypeFull type of episode is 'full'
	EpisodeTypeFull = "Full"
	// EpisodeTypeTrailer type of episode is 'trailer'
	EpisodeTypeTrailer = "Trailer"
	// EpisodeTypeBonus type of episode is 'bonus'
	EpisodeTypeBonus = "Bonus"

	// ResourceTypeExternal references an external URL
	ResourceTypeExternal = "external"
	// ResourceTypeLocal references a local resource
	ResourceTypeLocal = "local"
	// ResourceTypeImport references an external resources that will be imported into the CDN
	ResourceTypeImport = "import"

	// ResourceShow is referencing a resource of type "show"
	ResourceShow = "show"
	// ResourceEpisode is referencing a resource of type "episode"
	ResourceEpisode = "episode"
	// ResourceAsset is referencing any media or binary resource e.g. .mp3 or .png
	ResourceAsset = "asset"
	// ResourceALL is a wildcard for any kind of resource
	ResourceALL = "all"
)

type (
	// Apple Podcast: https://help.apple.com/itc/podcasts_connect/#/itcb54353390
	// RSS 2.0: https://cyber.harvard.edu/rss/rss.html

	// Metadata contains information describing a resource
	Metadata struct {
		Name   string            `json:"name" yaml:"name" binding:"required"` // REQUIRED <unique name>
		Labels map[string]string `json:"labels" yaml:"labels,omitempty"`      // REQUIRED
	}

	// GenericResource holds only the kind and metadata of a resource
	GenericResource struct {
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
		Image       Asset           `json:"image" yaml:"image" binding:"required"`             // REQUIRED 'channel.itunes.image'
	}

	// Episode holds all metadata related to a podcast episode
	Episode struct {
		APIVersion  string             `json:"apiVersion" yaml:"apiVersion" binding:"required"`   // REQUIRED default: v1.0
		Kind        string             `json:"kind" yaml:"kind" binding:"required"`               // REQUIRED default: episode
		Metadata    Metadata           `json:"metadata" yaml:"metadata" binding:"required"`       // REQUIRED
		Description EpisodeDescription `json:"description" yaml:"description" binding:"required"` // REQUIRED
		Image       Asset              `json:"image" yaml:"image" binding:"required"`             // REQUIRED 'item.itunes.image'
		Enclosure   Asset              `json:"enclosure" yaml:"enclosure" binding:"required"`     // REQUIRED
	}

	// ShowDescription holds essential show metadata
	ShowDescription struct {
		Title     string   `json:"title" yaml:"title" binding:"required"`          // REQUIRED 'channel.title' 'channel.itunes.title'
		Summary   string   `json:"summary" yaml:"summary" binding:"required"`      // REQUIRED 'channel.description'
		Link      Asset    `json:"link" yaml:"link"`                               // RECOMMENDED 'channel.link'
		Category  Category `json:"category" yaml:"category" binding:"required"`    // REQUIRED channel.category
		Owner     Owner    `json:"owner" yaml:"owner"`                             // RECOMMENDED 'channel.itunes.owner'
		Author    string   `json:"author" yaml:"author"`                           // RECOMMENDED 'channel.itunes.author'
		Copyright string   `json:"copyright,omitempty" yaml:"copyright,omitempty"` // OPTIONAL 'channel.copyright'
		NewFeed   *Asset   `json:"newFeed,omitempty" yaml:"newFeed,omitempty"`     // OPTIONAL channel.itunes.new-feed-url -> move to label
	}

	// EpisodeDescription holds essential episode metadata
	EpisodeDescription struct {
		Title       string `json:"title" yaml:"title" binding:"required"`                                 // REQUIRED 'item.title' 'item.itunes.title'
		Summary     string `json:"summary" yaml:"summary" binding:"required"`                             // REQUIRED 'item.description'
		EpisodeText string `json:"episodeText,omitempty" yaml:"episodeText,omitempty" binding:"required"` // REQUIRED 'item.itunes.summary'
		Link        Asset  `json:"link" yaml:"link"`                                                      // RECOMMENDED 'item.link'
		Duration    int    `json:"duration" yaml:"duration" binding:"required"`                           // REQUIRED 'item.itunes.duration'
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

	// Asset provides a link to a media resource
	Asset struct {
		URI   string `json:"uri" yaml:"uri" binding:"required"`      // REQUIRED
		Title string `json:"title,omitempty" yaml:"title,omitempty"` // OPTIONAL
		Rel   string `json:"rel,omitempty" yaml:"rel,omitempty"`     // OPTIONAL
		Type  string `json:"type,omitempty" yaml:"type,omitempty"`   // OPTIONAL
		Size  int    `json:"size,omitempty" yaml:"size,omitempty"`   // OPTIONAL
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

// ParentGUID is a convenience method to access the resources parent guid
func (e *Episode) Parent() string {
	return e.Metadata.Labels[LabelParentGUID]
}

// GUID is a convenience method to access the resources guid
func (r *GenericResource) GUID() string {
	return r.Metadata.Labels[LabelGUID]
}

// GUID is a convenience method to access the resources guid
func (s *Show) GUID() string {
	return s.Metadata.Labels[LabelGUID]
}

// ResolveURI re-writes the URI
func (r *Asset) ResolveURI(cdn, parent string) string {

	if r.Rel == ResourceTypeLocal {
		return fmt.Sprintf("%s/%s/%s", cdn, parent, r.URI)
	}
	if r.Rel == ResourceTypeImport {
		return fmt.Sprintf("%s/%s", cdn, metadata.FingerprintWithExt(parent, r.URI))
	}

	// r.Rel == ResourceTypeExternal or anything else, just return the URI as is ...
	return r.URI
}

func (r *Asset) AssetName() string {
	parts := strings.Split(r.URI, "/")
	if len(parts) == 0 {
		return r.URI
	}
	return parts[len(parts)-1]
}

// DefaultShowMetadata creates a default set of labels etc for a Show resource
//
//	language:	<ISO639 two-letter-code> REQUIRED 'channel.language'
//	explicit:	True | False REQUIRED 'channel.itunes.explicit'
//	type:		Episodic | Serial REQUIRED 'channel. itunes.type'
//	block:		Yes OPTIONAL 'channel.itunes.block' Anything else than 'Yes' has no effect
//	complete:	Yes OPTIONAL 'channel.itunes.complete' Anything else than 'Yes' has no effect
func DefaultShowMetadata(guid string) map[string]string {
	l := make(map[string]string)

	l[LabelLanguage] = "en_US"
	l[LabelExplicit] = "no"
	l[LabelType] = ShowTypeEpisodic
	l[LabelBlock] = "no"
	l[LabelComplete] = "no"
	l[LabelGUID] = guid

	return l
}

// DefaultEpisodeMetadata creates a default set of labels etc for a Episode resource
//	guid:		<unique id> 'item.guid'
//	date:		<publish date> REQUIRED 'item.pubDate'
//	season: 	<season number> OPTIONAL 'item.itunes.season'
//	episode:	<episode number> REQUIRED 'item.itunes.episode'
//	explicit:	True | False REQUIRED 'channel.itunes.explicit'
//	type:		Full | Trailer | Bonus REQUIRED 'item.itunes.episodeType'
//	block:		Yes OPTIONAL 'item.itunes.block' Anything else than 'Yes' has no effect
func DefaultEpisodeMetadata(guid, parent string) map[string]string {
	l := make(map[string]string)

	l[LabelGUID] = guid
	l[LabelParentGUID] = parent
	l[LabelDate] = time.Now().UTC().Format(time.RFC1123Z)
	l[LabelSeason] = "1"
	l[LabelEpisode] = "1"
	l[LabelExplicit] = "no"
	l[LabelType] = EpisodeTypeFull
	l[LabelBlock] = "no"

	return l
}

// DefaultShow creates a default show struc
func DefaultShow(name, title, summary, guid, portal, cdn string) *Show {
	return &Show{
		APIVersion: Version,
		Kind:       ResourceShow,
		Metadata: Metadata{
			Name:   name,
			Labels: DefaultShowMetadata(guid),
		},
		Description: ShowDescription{
			Title:   title,
			Summary: summary,
			Link: Asset{
				URI: fmt.Sprintf("%s/s/%s", portal, name),
			},
			Category: Category{
				Name: "Technology",
				SubCategory: []string{
					"Podcasting",
				},
			},
			Owner: Owner{
				Name:  fmt.Sprintf("%s owner", name),
				Email: fmt.Sprintf("hello@%s.me", name),
			},
			Author:    fmt.Sprintf("%s author", name),
			Copyright: fmt.Sprintf("%s copyright", name),
		},
		Image: Asset{
			URI: fmt.Sprintf("%s/default/cover.png", cdn),
			Rel: "external",
		},
	}
}

// DefaultEpisode creates a default episode struc
func DefaultEpisode(name, parentName, guid, parent, portal, cdn string) *Episode {
	return &Episode{
		APIVersion: Version,
		Kind:       ResourceEpisode,
		Metadata: Metadata{
			Name:   name,
			Labels: DefaultEpisodeMetadata(guid, parent),
		},
		Description: EpisodeDescription{
			Title:       fmt.Sprintf("%s - Episode Title", name),
			Summary:     fmt.Sprintf("%s - Episode Subtitle or short summary", name),
			EpisodeText: "A long-form description of the episode with notes etc.",
			Link: Asset{
				URI: fmt.Sprintf("%s/s/%s/%s", portal, parentName, name),
			},
			Duration: 1, // Seconds. Must not be 0, otherwise a validation error occurs.
		},
		Image: Asset{
			URI: fmt.Sprintf("%s/default/episode.png", cdn),
			Rel: "external",
		},
		Enclosure: Asset{
			URI:  fmt.Sprintf("%s/%s.mp3", parent, name),
			Type: "audio/mpeg",
			Rel:  "local",
		},
	}
}

func CreateGUID() string {
	guid, _ := id.ShortUUID()
	return guid
}
