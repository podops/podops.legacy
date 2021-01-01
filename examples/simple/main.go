package main

import (
	"fmt"
	"log"

	"github.com/podops/podops"
	a "github.com/podops/podops/apiv1"
	"github.com/txsvc/commons/pkg/util"
)

const (
	podcastName    string = "simple"
	podcastTitle   string = "Simple Podcast Example"
	podcastSummary string = "A podcast created using the PodOps API"

	defaultServiceEndpoint string = "http://localhost:8080"
	baseURL                string = "https://podops.dev"
)

func main() {

	// Initialize the client. This assumes that there is a configuration with
	// a valid token at 'DefaultConfigLocation' i.e. $HOME/.po/config
	client, err := podops.NewClientFromFile(podops.DefaultConfigLocation())
	if err != nil {
		log.Fatal(err)
	}

	// In order to test against a running local service,
	// set the API endpoint, otherwise https://api.podops.dev is assumes
	client.ServiceEndpoint = defaultServiceEndpoint

	// This verifies that the token is valid and we are allowed to access the API
	if err := client.Validate(); err != nil {
		log.Fatal(err)
	}

	// Create a new show/podcast. If a show 'podcastName' already exists (and we own it),
	// the call simply returns same metadata.
	p, err := client.CreateProduction(podcastName, podcastTitle, podcastSummary)
	if err != nil {
		log.Fatal(err)
	}

	// Set the context of the next operations to this show
	client.SetProduction(p.GUID)

	// Return a fully populated show struct with all defaults. We can add or
	// change values if needed.
	show := a.DefaultShow(baseURL, p.Name, podcastTitle, podcastSummary, p.GUID)
	show.Metadata.Labels[a.LabelLanguage] = "de_DE"
	show.Description.Author = "podops sample code"
	show.Description.Copyright = "Copyright 2021 podops.dev"
	show.Description.Owner.Name = "podops"
	show.Description.Owner.Email = "hello@podops.dev"

	// Each podcast needs at least one episode to be valid. Let's add one ...
	guid, _ := util.ShortUUID()
	episode := a.DefaultEpisode(baseURL, "first", show.Metadata.Name, guid, show.GUID())
	episode.Description.Title = "Drums !"
	episode.Description.Summary = "A short sample generated from one of Garage Band's default settings"
	episode.Description.EpisodeText = `This is the real description of the episode's content. Knock yourself out! 
	
	Withing 4000 characters that is ..`

	// Set the media data. This is a sample mp3 ...
	episode.Description.Duration = 21 // the duration of the episode, in seconds
	episode.Enclosure.URI = "https://cdn.podops.dev/default/sample.mp3"
	episode.Enclosure.Size = 503140 // bytes
	episode.Enclosure.Rel = "external"

	// Push changes to the show and episode to the service
	if _, err := client.UpdateResource(show.Kind, show.GUID(), true, &show); err != nil {
		log.Fatal(err)
	}
	if _, err := client.UpdateResource(episode.Kind, episode.GUID(), true, &episode); err != nil {
		log.Fatal(err)
	}

	// Start the feed build. On success the call returns the URL of the feed.xml
	feed, err := client.Build(show.GUID())
	if err != nil {
		log.Fatal(err)
	}

	// The path to the feed
	fmt.Printf("Access the podcast feed at %s\n", feed)
}
