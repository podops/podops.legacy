package metadata

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

	l[LabelLanguage] = "en"
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
func DefaultEpisodeMetadata(guid, parentGUID string) map[string]string {

	l := make(map[string]string)

	l[LabelGUID] = guid
	l[LabelParentGUID] = parentGUID
	l[LabelDate] = time.Now().UTC().Format(time.RFC1123Z)
	l[LabelSeason] = "1"
	l[LabelEpisode] = "1"
	l[LabelExplicit] = "no"
	l[LabelType] = EpisodeTypeFull
	l[LabelBlock] = "no"

	return l
}

// DefaultShow creates a default show struc
func DefaultShow(name, title, summary, guid string) *Show {
	return &Show{
		APIVersion: "v1",
		Kind:       "show",
		Metadata: Metadata{
			Name:   name,
			Labels: DefaultShowMetadata(guid),
		},
		Description: ShowDescription{
			Title:   title,
			Summary: summary,
			Link: Resource{
				URI: fmt.Sprintf("%s/s/%s", DefaultPortalEndpoint, name),
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
		Image: Resource{
			URI: fmt.Sprintf("%s/%s/coverart.png", DefaultCDNEndpoint, name),
			Rel: "external",
		},
	}
}

// DefaultEpisode creates a default episode struc
func DefaultEpisode(name, parentName, guid, parentGUID string) *Episode {

	return &Episode{
		APIVersion: "v1",
		Kind:       "episode",
		Metadata: Metadata{
			Name:   name,
			Labels: DefaultEpisodeMetadata(guid, parentGUID),
		},
		Description: EpisodeDescription{
			Title:       fmt.Sprintf("%s - Episode Title", name),
			Summary:     fmt.Sprintf("%s - Episode Subtitle or short summary", name),
			EpisodeText: "A long-form description of the episode with notes etc.",
			Link: Resource{
				URI: fmt.Sprintf("%s/s/%s/%s", DefaultPortalEndpoint, parentName, name),
			},
			Duration: 1, // Seconds. Must not be 0, otherwise a validation error occurs.
		},
		Image: Resource{
			URI: fmt.Sprintf("%s/%s/%s/coverart.png", DefaultCDNEndpoint, parentName, guid),
			Rel: "external",
		},
		Enclosure: Resource{
			URI:  fmt.Sprintf("%s/%s/%s/%s.mp3", DefaultCDNEndpoint, parentName, guid, name),
			Type: "audio/mpeg",
			Rel:  "external",
			Size: 0, // bytes
		},
	}
}
