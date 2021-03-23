package backend

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"cloud.google.com/go/storage"
	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/platform"
	a "github.com/podops/podops/apiv1"
	p "github.com/podops/podops/internal/platform"
	"google.golang.org/api/iterator"
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

// Build gathers all resources and builds the feed
func Build(ctx context.Context, guid string, validateOnly bool) error {

	var episodes EpisodeList

	p, err := GetProduction(ctx, guid)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("can not find '%s'", guid)
	}

	if err = ValidateProduction(ctx, guid); err != nil {
		p, err := GetProduction(ctx, guid)
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
	s, kind, _, err := ReadResource(ctx, fmt.Sprintf("%s/show-%s.yaml", guid, guid))
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
	obj := platform.Storage().Bucket(a.BucketCDN).Object(fmt.Sprintf("%s/feed.xml", guid))
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
func EnsureAsset(ctx context.Context, parent string, rsrc *a.Asset) error {
	if rsrc.Rel == a.ResourceTypeExternal {
		_, err := pingURL(rsrc.URI)
		return err
	}
	if rsrc.Rel == a.ResourceTypeLocal {
		path := fmt.Sprintf("%s/%s", parent, rsrc.URI)
		if !resourceExists(ctx, path) {
			return fmt.Errorf("can not find '%s'", rsrc.URI)
		}
		return nil
	}
	if rsrc.Rel == a.ResourceTypeImport {
		_, err := pingURL(rsrc.URI) // ping the URL already here to avoid queueing a request that will fail later anyways
		if err != nil {
			return err
		}

		path := rsrc.FingerprintURI(parent)
		if resourceExists(ctx, path) { // do nothing as the asset is present FIXME re-download if --force is set
			return nil // FIXME verify that the asset is unchanged, otherwise re-import
		}

		// dispatch a request for background import
		_, err = p.CreateTask(ctx, importTaskWithPrefix, &a.Import{Source: rsrc.URI, Dest: path})
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
	req.Header.Set("User-Agent", a.UserAgentString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		defer resp.Body.Close()
		// anything other than OK, Created, Accepted, NoContent is treated as an error
		if resp.StatusCode > http.StatusNoContent {
			return nil, fmt.Errorf("can not verify '%s'", url)
		}
	}
	return resp.Header.Clone(), nil
}

// resourceExists verifies the resource .yaml exists
func resourceExists(ctx context.Context, path string) bool {
	obj := platform.Storage().Bucket(a.BucketCDN).Object(path)
	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false
	}
	return true
}
