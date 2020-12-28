package resources

import (
	"context"
	"fmt"
	"sort"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/txsvc/platform/pkg/platform"

	t "github.com/podops/podops/internal/types"
	"github.com/podops/podops/pkg/metadata"
)

type (
	// EpisodeList holds the list of valid episodes that will be added to a podcast
	EpisodeList []*metadata.Episode
)

func (e EpisodeList) Len() int      { return len(e) }
func (e EpisodeList) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e EpisodeList) Less(i, j int) bool {
	return e[i].PublishDateTimestamp() > e[j].PublishDateTimestamp() // sorting direction is descending
}

// Build gathers all resources and builds the feed
func Build(ctx context.Context, guid string) error {

	var episodes EpisodeList

	p, err := GetProduction(ctx, guid)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("Can not find '%s'", guid)
	}

	// find all episodes and sort them by pubDate
	q := &storage.Query{
		Prefix: fmt.Sprintf("%s/episode", p.GUID),
	}
	it := platform.Storage().Bucket(t.BucketProduction).Objects(ctx, q)
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
		// FIXME skip episodes if block == yes etc
		episodes = append(episodes, e.(*metadata.Episode))
	}
	if episodes.Len() == 0 {
		return fmt.Errorf("Can not build feed with zero episodes")
	}

	sort.Sort(episodes)

	// read the show
	s, kind, _, err := ReadResource(ctx, fmt.Sprintf("%s/show-%s.yaml", guid, guid))
	if err != nil {
		return err
	}
	if kind != "show" {
		return fmt.Errorf("Unsupported resource '%s'", kind)
	}

	// build the feed XML
	show := s.(*metadata.Show)
	feed, err := show.Podcast()
	if err != nil {
		return err
	}

	tt, _ := time.Parse(time.RFC1123Z, episodes[0].PublishDate())
	feed.AddPubDate(&tt)

	for _, e := range episodes {
		item, err := e.Item()
		if err != nil {
			return err
		}
		feed.AddItem(item)
	}

	// dump the feed to the CDN location
	obj := platform.Storage().Bucket(t.BucketCDN).Object(fmt.Sprintf("%s/feed.xml", guid))
	writer := obj.NewWriter(ctx)
	if _, err := writer.Write(feed.Bytes()); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}
