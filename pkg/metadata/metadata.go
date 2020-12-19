package metadata

import (
	"time"
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

	LabelLanguage = "language"
	LabelExplicit = "explicit"
	LabelType     = "type"
	LabelBlock    = "block"
	LabelComplete = "complete"
	LabelGUID     = "guid"
	LabelDate     = "date"
	LabelSeason   = "season"
	LabelEpisode  = "episode"

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
		Name   string            // <unique name> REQUIRED
		Labels map[string]string `yaml:"labels,omitempty"` // REQUIRED
	}

	// Show holds all metadata related to a podcast/show
	Show struct {
		APIVersion  string          `yaml:"apiVersion"` // default: v1.0
		Kind        string          // default: show
		Metadata    Metadata        // REQUIRED
		Description ShowDescription // REQUIRED
		Image       Resource        // REQUIRED 'channel.itunes.image'
	}

	// Episode holds all metadata related to a podcast episode
	Episode struct {
		APIVersion  string             `yaml:"apiVersion"` // default: v1.0
		Kind        string             // default: episode
		Metadata    Metadata           // REQUIRED
		Description EpisodeDescription // REQUIRED
		Image       Resource           // REQUIRED 'item.itunes.image'
		Enclosure   Resource           // REQUIRED
	}

	// ShowDescription holds essential show metadata
	ShowDescription struct {
		Title     string    // REQUIRED 'channel.title' 'channel.itunes.title'
		Summary   string    // REQUIRED 'channel.description'
		Link      Resource  // RECOMMENDED 'channel.link'
		Category  Category  // REQUIRED channel.category
		Owner     Owner     // RECOMMENDED 'channel.itunes.owner'
		Author    string    // RECOMMENDED 'channel.itunes.author'
		Copyright string    `yaml:"copyright,omitempty"` // OPTIONAL 'channel.copyright'
		NewFeed   *Resource `yaml:"newFeed,omitempty"`   // OPTIONAL channel.itunes.new-feed-url -> move to label
	}

	// EpisodeDescription holds essential episode metadata
	EpisodeDescription struct {
		Title       string   // REQUIRED 'item.title' 'item.itunes.title'
		Summary     string   // REQUIRED 'item.description'
		EpisodeText string   // REQUIRED 'item.itunes.summary'
		Link        Resource // RECOMMENDED 'item.link'
		Duration    int      // REQUIRED 'item.itunes.duration'
	}

	// Owner describes the owner of the show/podcast
	Owner struct {
		Name  string // REQUIRED
		Email string // REQUIRED
	}

	// Category is the show/episodes category and it's subcategories
	Category struct {
		Name        string   // REQUIRED
		SubCategory []string `yaml:"subcategory,omitempty"` // OPTIONAL
	}

	// Resource provides a link to a media resource
	Resource struct {
		URI    string // REQUIRED
		Title  string `yaml:"title,omitempty"`
		Anchor string `yaml:"anchor,omitempty"`
		Rel    string `yaml:"rel,omitempty"`
		Type   string `yaml:"type,omitempty"`
		Size   int    `yaml:"size,omitempty"`
	}
)

// DefaultShowMetadata creates a default set of labels etc for a Show resource
//
//	language:	<ISO639 two-letter-code> REQUIRED 'channel.language'
//	explicit:	True | False REQUIRED 'channel.itunes.explicit'
//	type:		Episodic | Serial REQUIRED 'channel. itunes.type'
//	block:		Yes OPTIONAL 'channel.itunes.block' Anything else than 'Yes' has no effect
//	complete:	Yes OPTIONAL 'channel.itunes.complete' Anything else than 'Yes' has no effect
func DefaultShowMetadata() map[string]string {

	l := make(map[string]string)

	l[LabelLanguage] = "en"
	l[LabelExplicit] = "no"
	l[LabelType] = ShowTypeEpisodic
	l[LabelBlock] = "no"
	l[LabelComplete] = "no"

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
func DefaultEpisodeMetadata() map[string]string {

	l := make(map[string]string)

	l[LabelGUID] = "GUID"
	l[LabelDate] = time.Now().UTC().Format(time.RFC1123Z)
	l[LabelSeason] = "1"
	l[LabelEpisode] = "1"
	l[LabelExplicit] = "no"
	l[LabelType] = EpisodeTypeFull
	l[LabelBlock] = "no"

	return l
}
