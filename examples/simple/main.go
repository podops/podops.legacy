package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fupas/commons/pkg/util"

	"github.com/podops/podops"
	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/pkg/cli"
)

const (
	podcastName    string = "simple-podcast"
	podcastTitle   string = "PodOps Simple SDK Example"
	podcastSummary string = "A simple podcast for testing and experimentation. Created with the PodOps API."
)

// In order to test against a running local service,
// set API_ENDPOINT, otherwise https://api.podops.dev is used.

func main() {

	fmt.Println("\nStart creating a new podcast")
	start := util.TimestampNano()

	// Initialize the client. This assumes that there is a configuration with a valid token
	opts := cli.LoadConfiguration()
	client, err := podops.NewClient(context.TODO(), opts.Token)
	if err != nil {
		log.Fatal(err)
	}

	// This verifies that the token is valid and the client is allowed to access the API endpoint
	if !client.Valid() {
		log.Fatal(err)
	}
	fmt.Printf("\n[%d ms] Client authenticated.\n", delta(start))

	// Create a new show (podcast). If a show 'podcastName' already exists
	// and is owned by the token owner, the call simply returns its metadata.
	p, err := client.CreateProduction(podcastName, podcastTitle, podcastSummary)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n[%d ms] Registered a podcast.GUID='%s'\n", delta(start), p.GUID)

	// Set the context of the next operations to this show
	client.SetProduction(p.GUID)

	fmt.Printf("\n[%d ms] Select the podcast.\n", delta(start))

	// Return a fully populated show struct with all defaults.
	// We can add or change values if needed.
	show := a.DefaultShow(p.Name, podcastTitle, podcastSummary, p.GUID, client.PortalEndpoint(), client.CDNEndpoint())
	show.Metadata.Labels[a.LabelLanguage] = "de_DE"
	show.Description.Author = "PodOps sample code"
	show.Description.Copyright = "Copyright 2021 - Transformative Services"
	show.Description.Owner.Name = "Transformative Services"
	show.Description.Owner.Email = "hello@txs.vc"

	if _, err := client.UpdateResource(client.DefaultProduction(), show.Kind, show.GUID(), true, &show); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n[%d ms] Updated the show.\n", delta(start))

	// Each podcast needs at least one episode to be valid. Let's add one ...
	guid, _ := util.ShortUUID()
	episode := a.DefaultEpisode("first", show.Metadata.Name, guid, show.GUID(), client.PortalEndpoint(), client.CDNEndpoint())
	episode.Description.Title = "Drums !"
	episode.Description.Summary = "A short sample generated from one of Garage Band's default settings"
	episode.Description.EpisodeText = `This is the real description of the episode's content. Knock yourself out! 
	
	Withing 4000 characters that is ..`

	// Add the media data. This links to a sample mp3 ...
	episode.Description.Duration = 21 // the duration of the episode, in seconds
	episode.Enclosure.URI = fmt.Sprintf("%s/c/default/sample.mp3", a.DefaultCDNEndpoint)
	episode.Enclosure.Size = 503140 // bytes
	episode.Enclosure.Rel = "external"

	if _, err := client.UpdateResource(client.DefaultProduction(), episode.Kind, episode.GUID(), true, &episode); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n[%d ms] Updated the episode.\n", delta(start))

	// Start the feed build. On success the call returns the URL of the feed.xml
	feed, err := client.Build(show.GUID())
	if err != nil {
		log.Fatal(err)
	}

	// The URL to the feed
	fmt.Printf("\n[%d ms] Access the podcast feed at %s\n", delta(start), feed.FeedURL)
}

func delta(start int64) int64 {
	return (util.TimestampNano() - start) / 1000000
}
