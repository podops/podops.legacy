package resources

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/platform/pkg/services"

	"github.com/podops/podops/internal/config"
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
func Build(ctx context.Context, guid string, validateOnly bool) error {

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
	it := platform.Storage().Bucket(config.BucketProduction).Objects(ctx, q)
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

	if validateOnly {
		fmt.Printf(feed.String())
		return nil
	}

	// dump the feed to the CDN location
	obj := platform.Storage().Bucket(config.BucketCDN).Object(fmt.Sprintf("%s/feed.xml", guid))
	writer := obj.NewWriter(ctx)
	if _, err := writer.Write(feed.Bytes()); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}

// EnsureAsset validates the existence of the asset and imports it if necessary
func EnsureAsset(ctx context.Context, guid string, a *metadata.Resource) error {
	if a.Rel == metadata.ResourceTypeExternal {
		_, err := pingURL(a.URI)
		return err
	}
	if a.Rel == metadata.ResourceTypeLocal {
		path := fmt.Sprintf("%s/%s", guid, a.URI)
		if !Exists(ctx, path) {
			return fmt.Errorf("Can not find '%s'", a.URI)
		}
		return nil
	}
	if a.Rel == metadata.ResourceTypeImport {
		_, err := pingURL(a.URI) // ping the URL already here to avoid queueing a request that will fail later anyways
		if err != nil {
			return err
		}

		path := a.FingerprintURI(guid)
		if Exists(ctx, path) { // do nothing as the asset is present
			return nil // FIXME verify that the asset is unchanged, otherwise re-import
		}

		// dispatch a request for background import
		_, err = services.CreateTask(ctx, importTaskWithPrefix, &t.ImportRequest{Source: a.URI, Dest: path})
		if err != nil {
			return err
		}
	}
	return nil
}

// pingURL tries a HEAD or GET request to verify that 'url' exists and is reachable
func pingURL(url string) (http.Header, error) {

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		defer resp.Body.Close()
		// anything other than OK, Created, Accepted, NoContent is treated as an error
		if resp.StatusCode > http.StatusNoContent {
			return nil, fmt.Errorf("Can not verify '%s'", url)
		}
	}
	return resp.Header.Clone(), nil
}
