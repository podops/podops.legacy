package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	m "github.com/podops/podops/pkg/metadata"
)

// TemplateCommand creates a resource template with all default values
func TemplateCommand(c *cli.Context) error {
	template := c.Args().First()
	if template != "show" && template != "episode" {
		fmt.Println(fmt.Sprintf("\nDon't know how to create '%s'", template))
		return nil
	}

	if template == "show" {
		show := m.Show{
			APIVersion: "v1",
			Kind:       "show",
			Metadata: m.Metadata{
				Name:   "podcast-name",
				Labels: m.DefaultShowMetadata(),
			},
			Description: m.ShowDescription{
				Title:   "Podcast Title",
				Summary: "Podcast summary describing the podcast",
				Link: m.Resource{
					URI: "https://podcast.fm/podcast-name",
				},
				Category: m.Category{
					Name: "Technology",
					SubCategory: []string{
						"Podcasting",
					},
				},
				Owner: m.Owner{
					Name:  "Podcast Owner Name",
					Email: "hello@podcast.me",
				},
				Author:    "Podcast author",
				Copyright: "Podcast copyright",
			},
			Image: m.Resource{
				URI: "https://podcast.fm/podcast-name/coverart.png",
				Rel: "external",
			},
		}

		doc, err := yaml.Marshal(&show)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		ioutil.WriteFile(fmt.Sprintf("show-%s.yaml", show.Metadata.Labels[m.LabelGUID]), doc, 0644)
		fmt.Printf("--- show dump:\n\n%s\n\n", string(doc))

	} else {

		episode := m.Episode{
			APIVersion: "v1",
			Kind:       "episode",
			Metadata: m.Metadata{
				Name:   "episode1",
				Labels: m.DefaultEpisodeMetadata(),
			},
			Description: m.EpisodeDescription{
				Title:       "Episode Title",
				Summary:     "Episode Subtitle or short summary",
				EpisodeText: "A long-form description of the episode with notes etc.",
				Link: m.Resource{
					URI: "https://podcast.fm/podcast-name/episode1",
				},
				Duration: 0,
			},
			Image: m.Resource{
				URI: "https://podcast.fm/podcast-name/episode1/episode-coverart.png",
				Rel: "external",
			},
			Enclosure: m.Resource{
				URI:  "podcast.fm/podcast-name/episode1/episode1.mp3",
				Type: "audio/mpeg",
				Rel:  "external",
				Size: 0,
			},
		}

		doc, err := yaml.Marshal(&episode)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		ioutil.WriteFile(fmt.Sprintf("episode-%s.yaml", episode.Metadata.Labels[m.LabelGUID]), doc, 0644)
		fmt.Printf("--- episode dump:\n\n%s\n\n", string(doc))

	}

	return nil
}
