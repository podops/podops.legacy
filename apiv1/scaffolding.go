package apiv1

import (
	"fmt"
	"time"
)

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
			Size: 1, // bytes
		},
	}
}
