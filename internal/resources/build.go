package resources

import (
	"context"
	"fmt"
	"sort"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"

	"github.com/podops/podops/pkg/metadata"
	"github.com/txsvc/platform/pkg/platform"
)

type (
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
	it := platform.Storage().Bucket(bucketProduction).Objects(ctx, q)
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

	t, _ := time.Parse(time.RFC1123Z, episodes[0].PublishDate())
	feed.AddPubDate(&t)

	for _, e := range episodes {
		item, err := e.Item()
		if err != nil {
			return err
		}
		feed.AddItem(item)
	}

	// dump the feed to the CDN location
	obj := platform.Storage().Bucket(bucketCDN).Object(fmt.Sprintf("%s/feed.xml", guid))
	writer := obj.NewWriter(ctx)
	if _, err := writer.Write(feed.Bytes()); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}
