package metadata

import (
	"errors"
	"time"

	"github.com/podops/podops/pkg/feed"
)

var mediaTypeMap map[string]feed.EnclosureType

func init() {
	mediaTypeMap = make(map[string]feed.EnclosureType)
	mediaTypeMap["audio/x-m4a"] = feed.M4A
	mediaTypeMap["video/x-m4v"] = feed.M4V
	mediaTypeMap["video/mp4"] = feed.MP4
	mediaTypeMap["audio/mpeg"] = feed.MP3
	mediaTypeMap["video/quicktime"] = feed.MOV
	mediaTypeMap["application/pdf"] = feed.PDF
	mediaTypeMap["document/x-epub"] = feed.EPUB
}

// Podcast transforms Show metadata into a podcast feed struct
func (s *Show) Podcast() (*feed.Podcast, error) {
	now := time.Now()

	// basics
	pf := feed.New(s.Description.Title, s.Description.Link.URI, s.Description.Summary, &now, &now)
	// details
	pf.AddSummary(s.Description.Summary)
	if s.Description.Author == "" {
		pf.AddAuthor(s.Description.Owner.Name, s.Description.Owner.Email)
	} else {
		pf.IAuthor = s.Description.Author
	}
	pf.AddCategory(s.Description.Category.Name, s.Description.Category.SubCategory)
	pf.AddImage(s.Image.URI)
	pf.IOwner = &feed.Author{
		Name:  s.Description.Owner.Name,
		Email: s.Description.Owner.Email,
	}
	pf.Copyright = s.Description.Copyright
	if s.Description.NewFeed != nil {
		pf.INewFeedURL = s.Description.NewFeed.URI
	}
	pf.Language = s.Metadata.Labels[LabelLanguage]
	pf.IExplicit = s.Metadata.Labels[LabelExplicit]

	t := s.Metadata.Labels[LabelType]
	if t == ShowTypeEpisodic || t == ShowTypeSerial {
		pf.IType = t
	} else {
		return nil, errors.New("Show type must be 'Episodic' or 'Serial' ")
	}
	if s.Metadata.Labels[LabelBlock] == "yes" {
		pf.IBlock = "yes"
	}
	if s.Metadata.Labels[LabelComplete] == "yes" {
		pf.IComplete = "yes"
	}

	return &pf, nil
}

//
//	explicit:	True | False REQUIRED 'channel.itunes.explicit'
//	type:		Episodic | Serial REQUIRED 'channel. itunes.type'
//	block:		Yes OPTIONAL 'channel.itunes.block' Anything else than 'Yes' has no effect
//	complete:	Yes OPTIONAL 'channel.itunes.complete' Anything else than 'Yes' has no effect

// Item returns the episode struct needed for a podcast feed struct
func (e *Episode) Item() (*feed.Item, error) {

	pubDate, err := time.Parse(time.RFC1123Z, e.Metadata.Labels[LabelDate])
	if err != nil {
		return nil, err
	}

	ef := &feed.Item{
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
	ef.GUID = e.Metadata.Labels[LabelGUID]
	ef.IExplicit = e.Metadata.Labels[LabelExplicit]
	ef.ISeason = e.Metadata.Labels[LabelSeason]
	ef.IEpisode = e.Metadata.Labels[LabelEpisode]
	ef.IEpisodeType = e.Metadata.Labels[LabelType]
	if e.Metadata.Labels[LabelBlock] == "yes" {
		ef.IBlock = "yes"
	}

	return ef, nil
}
