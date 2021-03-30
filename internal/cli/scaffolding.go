package cli

import (
	"fmt"
	"time"

	a "github.com/podops/podops"
)

// FIXME move this to pkg/cli

// DefaultShowMetadata creates a default set of labels etc for a Show resource
//
//	language:	<ISO639 two-letter-code> REQUIRED 'channel.language'
//	explicit:	True | False REQUIRED 'channel.itunes.explicit'
//	type:		Episodic | Serial REQUIRED 'channel. itunes.type'
//	block:		Yes OPTIONAL 'channel.itunes.block' Anything else than 'Yes' has no effect
//	complete:	Yes OPTIONAL 'channel.itunes.complete' Anything else than 'Yes' has no effect
func DefaultShowMetadata(guid string) map[string]string {
	l := make(map[string]string)

	l[a.LabelLanguage] = "en_US"
	l[a.LabelExplicit] = "no"
	l[a.LabelType] = a.ShowTypeEpisodic
	l[a.LabelBlock] = "no"
	l[a.LabelComplete] = "no"
	l[a.LabelGUID] = guid

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

	l[a.LabelGUID] = guid
	l[a.LabelParentGUID] = parent
	l[a.LabelDate] = time.Now().UTC().Format(time.RFC1123Z)
	l[a.LabelSeason] = "1"
	l[a.LabelEpisode] = "1"
	l[a.LabelExplicit] = "no"
	l[a.LabelType] = a.EpisodeTypeFull
	l[a.LabelBlock] = "no"

	return l
}

// DefaultShow creates a default show struc
func DefaultShow(name, title, summary, guid, portal, cdn string) *a.Show {
	return &a.Show{
		APIVersion: a.Version,
		Kind:       a.ResourceShow,
		Metadata: a.Metadata{
			Name:   name,
			Labels: DefaultShowMetadata(guid),
		},
		Description: a.ShowDescription{
			Title:   title,
			Summary: summary,
			Link: a.Asset{
				URI: fmt.Sprintf("%s/s/%s", portal, name),
			},
			Category: a.Category{
				Name: "Technology",
				SubCategory: []string{
					"Podcasting",
				},
			},
			Owner: a.Owner{
				Name:  fmt.Sprintf("%s owner", name),
				Email: fmt.Sprintf("hello@%s.me", name),
			},
			Author:    fmt.Sprintf("%s author", name),
			Copyright: fmt.Sprintf("%s copyright", name),
		},
		Image: a.Asset{
			URI: fmt.Sprintf("%s/c/default/cover.png", cdn), // FIXME create a real resource entry to be consistent
			Rel: "external",
		},
	}
}

// DefaultEpisode creates a default episode struc
func DefaultEpisode(name, parentName, guid, parent, portal, cdn string) *a.Episode {
	return &a.Episode{
		APIVersion: a.Version,
		Kind:       a.ResourceEpisode,
		Metadata: a.Metadata{
			Name:   name,
			Labels: DefaultEpisodeMetadata(guid, parent),
		},
		Description: a.EpisodeDescription{
			Title:       fmt.Sprintf("%s - Episode Title", name),
			Summary:     fmt.Sprintf("%s - Episode Subtitle or short summary", name),
			EpisodeText: "A long-form description of the episode with notes etc.",
			Link: a.Asset{
				URI: fmt.Sprintf("%s/s/%s/%s", portal, parentName, name),
			},
			Duration: 1, // Seconds. Must not be 0, otherwise a validation error occurs.
		},
		Image: a.Asset{
			URI: fmt.Sprintf("%s/c/default/episode.png", cdn), // FIXME create a real resource entry to be consistent
			Rel: "external",
		},
		Enclosure: a.Asset{
			URI:  fmt.Sprintf("%s/%s.mp3", parent, name),
			Type: "audio/mpeg",
			Rel:  "local",
			Size: 1, // bytes
		},
	}
}
