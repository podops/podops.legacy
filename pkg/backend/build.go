package backend

import (
	"context"
	"fmt"
	"sort"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/platform"

	a "github.com/podops/podops/apiv1"
)

type (
	// EpisodeList holds the list of valid episodes that will be added to a podcast
	EpisodeList []*a.Episode
)

func (e EpisodeList) Len() int      { return len(e) }
func (e EpisodeList) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e EpisodeList) Less(i, j int) bool {
	return e[i].PublishDateTimestamp() > e[j].PublishDateTimestamp() // sorting direction is descending
}

// BuildFeed gathers all resources and builds the feed
func BuildFeed(ctx context.Context, production string, validateOnly bool) error {

	var episodes EpisodeList

	p, err := GetProduction(ctx, production)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("can not find '%s'", production)
	}

	if err = ValidateProduction(ctx, production); err != nil {
		p, err := GetProduction(ctx, production)
		if err != nil {
			return err
		}
		p.BuildDate = 0 // FIXME BuildDate is the only flag we currently have to mark a production as VALID
		UpdateProduction(ctx, p)

		return fmt.Errorf("can not build feed")
	}

	// FIXME build this new !

	// find all episodes and sort them by pubDate
	q := &storage.Query{
		Prefix: fmt.Sprintf("%s/episode", p.GUID),
	}
	it := platform.Storage().Bucket(a.BucketProduction).Objects(ctx, q)
	for {
		attr, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		e, _, _, err := ReadResource(ctx, attr.Name)
		if err != nil {
			return err
		}
		episode := e.(*a.Episode)

		// FIXME skip episodes if block == yes or publish date is in the future
		if episode.PublishDateTimestamp() > util.Timestamp() {
			continue
		}
		if episode.Metadata.Labels[a.LabelBlock] == "yes" {
			continue
		}

		episodes = append(episodes, episode)
	}
	if episodes.Len() == 0 {
		return fmt.Errorf("can not build feed with zero episodes")
	}

	sort.Sort(episodes)

	// read the show
	s, kind, _, err := ReadResource(ctx, fmt.Sprintf("%s/show-%s.yaml", production, production))
	if err != nil {
		return err
	}
	if kind != a.ResourceShow {
		return fmt.Errorf("unsupported resource '%s'", kind)
	}

	// build the feed XML
	show := s.(*a.Show)
	feed, err := a.TransformToPodcast(show)
	if err != nil {
		return err
	}

	tt, _ := time.Parse(time.RFC1123Z, episodes[0].PublishDate())
	feed.AddPubDate(&tt)

	// FIXME use a -f flag to enforce asset assurance on build

	for _, e := range episodes {
		item, err := a.TransformToItem(e)
		if err != nil {
			return err
		}
		feed.AddItem(item)
	}

	if validateOnly {
		fmt.Printf(feed.String())
		return nil
	}

	// dump the feed to the CDN location
	obj := platform.Storage().Bucket(a.BucketCDN).Object(fmt.Sprintf("%s/feed.xml", production))
	writer := obj.NewWriter(ctx)
	if _, err := writer.Write(feed.Bytes()); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}
