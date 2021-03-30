package feed

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/platform"

	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/feed/rss"
	"github.com/podops/podops/pkg/backend"
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

var mediaTypeMap map[string]rss.EnclosureType

func init() {
	mediaTypeMap = make(map[string]rss.EnclosureType)
	mediaTypeMap["audio/x-m4a"] = rss.M4A
	mediaTypeMap["video/x-m4v"] = rss.M4V
	mediaTypeMap["video/mp4"] = rss.MP4
	mediaTypeMap["audio/mpeg"] = rss.MP3
	mediaTypeMap["video/quicktime"] = rss.MOV
	mediaTypeMap["application/pdf"] = rss.PDF
	mediaTypeMap["document/x-epub"] = rss.EPUB
}

// Build gathers all resources and builds the feed.xml
func Build(ctx context.Context, production string, validateOnly bool) error {

	var episodes EpisodeList

	p, err := backend.GetProduction(ctx, production)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("can not find '%s'", production)
	}

	if err = backend.ValidateProduction(ctx, production); err != nil {
		p, err := backend.GetProduction(ctx, production)
		if err != nil {
			return err
		}
		p.BuildDate = 0 // FIXME BuildDate is the only flag we currently have to mark a production as VALID
		backend.UpdateProduction(ctx, p)

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

		e, _, _, err := backend.ReadResource(ctx, attr.Name)
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
	s, kind, _, err := backend.ReadResource(ctx, fmt.Sprintf("%s/show-%s.yaml", production, production))
	if err != nil {
		return err
	}
	if kind != a.ResourceShow {
		return fmt.Errorf("unsupported resource '%s'", kind)
	}

	// build the feed XML
	show := s.(*a.Show)
	feed, err := TransformToPodcast(show)
	if err != nil {
		return err
	}

	tt, _ := time.Parse(time.RFC1123Z, episodes[0].PublishDate())
	feed.AddPubDate(&tt)

	// FIXME use a -f flag to enforce asset assurance on build

	for _, e := range episodes {
		item, err := TransformToItem(e)
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

// TransformToPodcast transforms Show metadata into a podcast feed struct
func TransformToPodcast(s *a.Show) (*rss.Channel, error) {
	now := time.Now()

	// basics
	pf := rss.New(s.Description.Title, s.Description.Link.URI, s.Description.Summary, &now, &now) // FIXME remove timestamps
	// details
	pf.AddSummary(s.Description.Summary)
	if s.Description.Author == "" {
		pf.AddAuthor(s.Description.Owner.Name, s.Description.Owner.Email)
	} else {
		pf.IAuthor = s.Description.Author
	}
	pf.AddCategory(s.Description.Category.Name, s.Description.Category.SubCategory)
	pf.AddImage(s.Image.ResolveURI(a.StorageEndpoint, s.GUID()))
	pf.IOwner = &rss.Author{
		Name:  s.Description.Owner.Name,
		Email: s.Description.Owner.Email,
	}
	pf.Copyright = s.Description.Copyright
	if s.Description.NewFeed != nil {
		pf.INewFeedURL = s.Description.NewFeed.URI
	}
	pf.Language = s.Metadata.Labels[a.LabelLanguage]
	pf.IExplicit = s.Metadata.Labels[a.LabelExplicit]

	t := s.Metadata.Labels[a.LabelType]
	if t == a.ShowTypeEpisodic || t == a.ShowTypeSerial {
		pf.IType = t
	} else {
		return nil, errors.New("Show type must be 'Episodic' or 'Serial' ")
	}
	if s.Metadata.Labels[a.LabelBlock] == "yes" {
		pf.IBlock = "yes"
	}
	if s.Metadata.Labels[a.LabelComplete] == "yes" {
		pf.IComplete = "yes"
	}

	return &pf, nil
}

//	explicit:	True | False REQUIRED 'channel.itunes.explicit'
//	type:		Episodic | Serial REQUIRED 'channel. itunes.type'
//	block:		Yes OPTIONAL 'channel.itunes.block' Anything else than 'Yes' has no effect
//	complete:	Yes OPTIONAL 'channel.itunes.complete' Anything else than 'Yes' has no effect

// TransformToItem returns the episode struct needed for a podcast feed struct
func TransformToItem(e *a.Episode) (*rss.Item, error) {

	pubDate, err := time.Parse(time.RFC1123Z, e.Metadata.Labels[a.LabelDate])
	if err != nil {
		return nil, err
	}

	ef := &rss.Item{
		Title:       e.Description.Title,
		Description: e.Description.Summary,
	}

	ef.AddEnclosure(e.Enclosure.ResolveURI(a.DefaultCDNEndpoint+"/c", e.Parent()), mediaTypeMap[e.Enclosure.Type], (int64)(e.Enclosure.Size))
	ef.AddImage(e.Image.ResolveURI(a.StorageEndpoint, e.Parent()))
	ef.AddPubDate(&pubDate)
	ef.AddSummary(e.Description.EpisodeText)
	ef.AddDuration((int64)(e.Description.Duration))
	ef.Link = e.Description.Link.URI
	ef.ISubtitle = e.Description.Summary
	ef.GUID = e.Metadata.Labels[a.LabelGUID]
	ef.IExplicit = e.Metadata.Labels[a.LabelExplicit]
	ef.ISeason = e.Metadata.Labels[a.LabelSeason]
	ef.IEpisode = e.Metadata.Labels[a.LabelEpisode]
	ef.IEpisodeType = e.Metadata.Labels[a.LabelType]
	if e.Metadata.Labels[a.LabelBlock] == "yes" {
		ef.IBlock = "yes"
	}

	return ef, nil
}
