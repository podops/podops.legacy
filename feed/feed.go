package feed

import (
	"context"
	"fmt"
	"time"

	"github.com/txsvc/platform/v2"
	ds "github.com/txsvc/platform/v2/pkg/datastore"
	"github.com/txsvc/platform/v2/pkg/timestamp"

	"github.com/podops/podops"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/feed/rss"
	"github.com/podops/podops/internal/errordef"
	"github.com/podops/podops/internal/messagedef"
)

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

	p, err := backend.GetProduction(ctx, production)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf(messagedef.MsgResourceNotFound, production)
	}

	if err = backend.ValidateProduction(ctx, production); err != nil {
		p, err := backend.GetProduction(ctx, production)
		if err != nil {
			return err
		}
		p.BuildDate = 0
		p.Published = false
		p.LatestPublishDate = 0

		backend.UpdateProduction(ctx, p)

		return errordef.ErrFeedFailed
	}

	// list all episodes, excluding future (i.e. unpublished) ones, descending order

	now := timestamp.Now()
	er, err := backend.ListPublishedEpisodes(ctx, production, now, 1)
	if err != nil {
		platform.ReportError(err)
		return err
	}

	if len(er) == 0 {
		return errordef.ErrFeedFailed
	}

	// read all episodes
	episodes := make([]*podops.Episode, len(er))
	for i := range er {
		e, err := backend.GetResourceContent(ctx, er[i].GUID)
		if err != nil {
			return err
		}
		// FIXME filter for other flags, e.g. Block = true
		episodes[i] = e.(*podops.Episode)
	}

	// read the show
	s, err := backend.GetResourceContent(ctx, production)
	if err != nil {
		return err
	}

	// build the feed XML
	show := s.(*podops.Show)
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
		return nil // no errors so far, the feed is valid
	}

	// dump the feed to the CDN
	obj := ds.Storage().Bucket(podops.BucketProduction).Object(fmt.Sprintf("%s/feed.xml", production))
	writer := obj.NewWriter(ctx)
	if _, err := writer.Write(feed.Bytes()); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	// source data should be OK by now, we can update the metadata
	p.BuildDate = timestamp.Now()
	p.Published = true
	p.LatestPublishDate = er[0].Published
	if err := backend.UpdateProduction(ctx, p); err != nil {
		return err
	}

	return nil
}

// TransformToPodcast transforms Show metadata into a podcast feed struct
func TransformToPodcast(s *podops.Show) (*rss.Channel, error) {
	now := time.Now()

	// basics
	pf := rss.New(s.Description.Title, s.Description.Link.URI, s.Description.Summary, &now, &now)
	// details
	pf.AddSummary(s.Description.Summary)
	if s.Description.Author == "" {
		pf.AddAuthor(s.Description.Owner.Name, s.Description.Owner.Email)
	} else {
		pf.IAuthor = s.Description.Author
	}
	pf.AddCategory(s.Description.Category.Name, s.Description.Category.SubCategory)
	pf.AddImage(s.Image.URI)
	pf.IOwner = &rss.Author{
		Name:  s.Description.Owner.Name,
		Email: s.Description.Owner.Email,
	}
	pf.Copyright = s.Description.Copyright
	if s.Description.NewFeed != nil {
		pf.INewFeedURL = s.Description.NewFeed.URI
	}
	pf.Language = s.Metadata.Labels[podops.LabelLanguage]
	pf.IExplicit = s.Metadata.Labels[podops.LabelExplicit]

	t := s.Metadata.Labels[podops.LabelType]
	if t == podops.ShowTypeEpisodic || t == podops.ShowTypeSerial {
		pf.IType = t
	} else {
		return nil, errordef.ErrInvalidParameters
	}
	if s.Metadata.Labels[podops.LabelBlock] == "yes" {
		pf.IBlock = "yes"
	}
	if s.Metadata.Labels[podops.LabelComplete] == "yes" {
		pf.IComplete = "yes"
	}

	return &pf, nil
}

//	explicit:	True | False REQUIRED 'channel.itunes.explicit'
//	type:		Episodic | Serial REQUIRED 'channel. itunes.type'
//	block:		Yes OPTIONAL 'channel.itunes.block' Anything else than 'Yes' has no effect
//	complete:	Yes OPTIONAL 'channel.itunes.complete' Anything else than 'Yes' has no effect

// TransformToItem returns the episode struct needed for a podcast feed struct
func TransformToItem(e *podops.Episode) (*rss.Item, error) {

	pubDate, err := time.Parse(time.RFC1123Z, e.Metadata.Labels[podops.LabelDate])
	if err != nil {
		return nil, err
	}

	ef := &rss.Item{
		Title:       e.Description.Title,
		Description: e.Description.Summary,
	}

	ef.AddEnclosure(e.Enclosure.URI, mediaTypeMap[e.Enclosure.Type], (int64)(e.Enclosure.Size))
	ef.AddImage(e.Image.URI)
	ef.AddPubDate(&pubDate)
	ef.AddSummary(e.Description.EpisodeText)
	ef.AddDuration((int64)(e.Description.Duration))
	ef.Link = e.Description.Link.URI
	ef.ISubtitle = e.Description.Summary
	ef.GUID = e.Metadata.Labels[podops.LabelGUID]
	ef.IExplicit = e.Metadata.Labels[podops.LabelExplicit]
	ef.ISeason = e.Metadata.Labels[podops.LabelSeason]
	ef.IEpisode = e.Metadata.Labels[podops.LabelEpisode]
	ef.IEpisodeType = e.Metadata.Labels[podops.LabelType]
	if e.Metadata.Labels[podops.LabelBlock] == "yes" {
		ef.IBlock = "yes"
	}

	return ef, nil
}
