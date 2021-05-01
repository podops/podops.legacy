package rss

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/errordef"
)

// iTunes Specifications: https://help.apple.com/itc/podcasts_connect/#/itcb54353390

const (
	M4A EnclosureType = iota
	M4V
	MP4
	MP3
	MOV
	PDF
	EPUB

	enclosureDefault = "application/octet-stream"
)

// New instantiates a podcast with required parameters.
//
// Nil-able fields are optional but recommended as they are formatted to the expected proper formats.
func New(title, link, description string, pubDate, lastBuildDate *time.Time) Channel {
	return Channel{
		Title:         title,
		ITitle:        title,
		Link:          link,
		Description:   description,
		Generator:     fmt.Sprintf("PodOps %s (https://github.com/podops/podops)", podops.VersionString),
		PubDate:       parseDateRFC1123Z(pubDate),
		LastBuildDate: parseDateRFC1123Z(lastBuildDate),
		Language:      "en",

		// setup dependency (could inject later)
		encode: encoder,
	}
}

// AddAuthor adds the specified Author to the podcast.
func (p *Channel) AddAuthor(name, email string) {
	if len(email) == 0 {
		return
	}
	p.ManagingEditor = parseAuthorNameEmail(&Author{
		Name:  name,
		Email: email,
	})
	p.IAuthor = p.ManagingEditor
}

// AddAtomLink adds a FQDN reference to an atom feed.
func (p *Channel) AddAtomLink(href string) {
	if len(href) == 0 {
		return
	}
	p.AtomLink = &AtomLink{
		HREF: href,
		Rel:  "self",
		Type: "application/rss+xml",
	}
}

// AddCategory adds the category to the podcast.
//
// ICategory can be listed multiple times.
//
// Calling this method multiple times will APPEND the category to the existing
// list, if any, including ICategory.
//
// Note that Apple iTunes has a specific list of categories that only can be
// used and will invalidate the feed if deviated from the list.
//
func (p *Channel) AddCategory(category string, subCategories []string) {
	if len(category) == 0 {
		return
	}

	// RSS 2.0 Category only supports 1-tier
	if len(p.Category) > 0 {
		p.Category = p.Category + "," + category
	} else {
		p.Category = category
	}

	icat := ICategory{Text: category}
	for _, c := range subCategories {
		if len(c) == 0 {
			continue
		}
		icat2 := ICategory{Text: c}
		icat.ICategories = append(icat.ICategories, &icat2)
	}
	p.ICategories = append(p.ICategories, &icat)
}

// AddImage adds the specified Image to the Podcast.
//
// Podcast feeds contain artwork that is a minimum size of
// 1400 x 1400 pixels and a maximum size of 3000 x 3000 pixels,
// 72 dpi, in JPEG or PNG format with appropriate file
// extensions (.jpg, .png), and in the RGB colorspace. To optimize
// images for mobile devices, Apple recommends compressing your
// image files.
func (p *Channel) AddImage(url string) {
	if len(url) == 0 {
		return
	}
	p.Image = &Image{
		URL:   url,
		Title: p.Title,
		Link:  p.Link,
	}
	p.IImage = &IImage{HREF: url}
}

// AddItem adds the podcast episode.  It returns a count of Items added or any
// errors in validation that may have occurred.
//
// This method takes the "itunes overrides" approach to populating
// itunes tags according to the overrides rules in the specification.
// This not only complies completely with iTunes parsing rules; but, it also
// displays what is possible to be set on an individual episode level – if you
// wish to have more fine grain control over your content.
//
// This method imposes strict validation of the Item being added to confirm
// to Podcast and iTunes specifications.
//
// Article minimal requirements are:
//
//   * Title
//   * Description
//   * Link
//
// Audio, Video and Downloads minimal requirements are:
//
//   * Title
//   * Description
//   * Enclosure (HREF, Type and Length all required)
//
// The following fields are always overwritten (don't set them):
//
//   * GUID
//   * PubDateFormatted
//   * AuthorFormatted
//   * Enclosure.TypeFormatted
//   * Enclosure.LengthFormatted
//
// Recommendations:
//
//   * Just set the minimal fields: the rest get set for you.
//   * Always set an Enclosure.Length, to be nice to your downloaders.
//   * Follow Apple's best practices to enrich your podcasts:
//     https://help.apple.com/itc/podcasts_connect/#/itc2b3780e76
//   * For specifications of itunes tags, see:
//     https://help.apple.com/itc/podcasts_connect/#/itcb54353390
//
func (p *Channel) AddItem(i *Item) (int, error) {
	// initial guards for required fields
	if len(i.Title) == 0 || len(i.Description) == 0 {
		return len(p.Items), errors.New("Title and Description are required")
	}
	if i.Enclosure != nil {
		if len(i.Enclosure.URL) == 0 {
			return len(p.Items),
				errors.New(i.Title + ": Enclosure.URL is required")
		}
		if i.Enclosure.Type.String() == enclosureDefault {
			return len(p.Items),
				errors.New(i.Title + ": Enclosure.Type is required")
		}
	} else if len(i.Link) == 0 {
		return len(p.Items),
			errors.New(i.Title + ": Link is required when not using Enclosure")
	}

	// corrective actions and overrides

	i.PubDateFormatted = parseDateRFC1123Z(i.PubDate)
	i.AuthorFormatted = parseAuthorNameEmail(i.Author)
	if i.Enclosure != nil {
		if len(i.GUID) == 0 {
			i.GUID = i.Enclosure.URL // yep, GUID is the Permlink URL
		}

		if i.Enclosure.Length < 0 {
			i.Enclosure.Length = 0
		}
		i.Enclosure.LengthFormatted = strconv.FormatInt(i.Enclosure.Length, 10)
		i.Enclosure.TypeFormatted = i.Enclosure.Type.String()

		// allow Link to be set for article references to Downloads,
		// otherwise set it to the enclosurer's URL.
		if len(i.Link) == 0 {
			i.Link = i.Enclosure.URL
		}
	} else {
		i.GUID = i.Link // yep, GUID is the Permlink URL
	}

	// iTunes it

	if len(i.IAuthor) == 0 {
		switch {
		case i.Author != nil:
			i.IAuthor = i.Author.Email
		case len(p.IAuthor) != 0:
			i.Author = &Author{Email: p.IAuthor}
			i.IAuthor = p.IAuthor
		case len(p.ManagingEditor) != 0:
			i.Author = &Author{Email: p.ManagingEditor}
			i.IAuthor = p.ManagingEditor
		}
	}
	if i.IImage == nil {
		if p.Image != nil {
			i.IImage = &IImage{HREF: p.Image.URL}
		}
	}

	p.Items = append(p.Items, i)
	return len(p.Items), nil
}

// AddPubDate adds the datetime as a parsed PubDate.
//
// UTC time is used by default.
func (p *Channel) AddPubDate(datetime *time.Time) {
	p.PubDate = parseDateRFC1123Z(datetime)
}

// AddLastBuildDate adds the datetime as a parsed PubDate.
//
// UTC time is used by default.
func (p *Channel) AddLastBuildDate(datetime *time.Time) {
	p.LastBuildDate = parseDateRFC1123Z(datetime)
}

// AddSubTitle adds the iTunes subtitle that is displayed with the title
// in iTunes.
//
// Note that this field should be just a few words long according to Apple.
// This method will truncate the string to 64 chars if too long with "..."
func (p *Channel) AddSubTitle(subTitle string) {
	count := utf8.RuneCountInString(subTitle)
	if count == 0 {
		return
	}
	if count > 64 {
		s := []rune(subTitle)
		subTitle = string(s[0:61]) + "..."
	}
	p.ISubtitle = subTitle
}

// AddSummary adds the iTunes summary.
//
// Limit: 4000 characters
//
// Note that this field is a CDATA encoded field which allows for rich text
// such as html links: `<a href="http://www.apple.com">Apple</a>`.
func (p *Channel) AddSummary(summary string) {
	count := utf8.RuneCountInString(summary)
	if count == 0 {
		return
	}
	if count > 4000 {
		s := []rune(summary)
		summary = string(s[0:4000])
	}
	p.ISummary = &ISummary{
		Text: summary,
	}
}

// Bytes returns an encoded []byte slice.
func (p *Channel) Bytes() []byte {
	return []byte(p.String())
}

// Encode writes the bytes to the io.Writer stream in RSS 2.0 specification.
func (p *Channel) Encode(w io.Writer) error {
	if _, err := w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")); err != nil {
		return errordef.ErrBuildFailed
	}

	atomLink := ""
	if p.AtomLink != nil {
		atomLink = "http://www.w3.org/2005/Atom"
	}
	wrapped := channelWrapper{
		ITUNESNS: "http://www.itunes.com/dtds/podcast-1.0.dtd",
		ATOMNS:   atomLink,
		Version:  "2.0",
		Channel:  p,
	}
	return p.encode(w, wrapped)
}

// String encodes the Podcast state to a string.
func (p *Channel) String() string {
	b := new(bytes.Buffer)
	if err := p.Encode(b); err != nil {
		return "String: channel.write returned the error: " + err.Error()
	}
	return b.String()
}

// AddEnclosure adds the downloadable asset to the podcast Item.
func (i *Item) AddEnclosure(
	url string, enclosureType EnclosureType, lengthInBytes int64) {
	i.Enclosure = &Enclosure{
		URL:    url,
		Type:   enclosureType,
		Length: lengthInBytes,
	}
}

// AddImage adds the image as an iTunes-only IImage.  RSS 2.0 does not have
// the specification of Images at the Item level.
//
// Podcast feeds contain artwork that is a minimum size of
// 1400 x 1400 pixels and a maximum size of 3000 x 3000 pixels,
// 72 dpi, in JPEG or PNG format with appropriate file
// extensions (.jpg, .png), and in the RGB colorspace. To optimize
// images for mobile devices, Apple recommends compressing your
// image files.
func (i *Item) AddImage(url string) {
	if len(url) > 0 {
		i.IImage = &IImage{HREF: url}
	}
}

// AddPubDate adds the datetime as a parsed PubDate.
//
// UTC time is used by default.
func (i *Item) AddPubDate(datetime *time.Time) {
	i.PubDate = datetime
	i.PubDateFormatted = parseDateRFC1123Z(i.PubDate)
}

// AddSummary adds the iTunes summary.
//
// Limit: 4000 characters
//
// Note that this field is a CDATA encoded field which allows for rich text
// such as html links: `<a href="http://www.apple.com">Apple</a>`.
func (i *Item) AddSummary(summary string) {
	count := utf8.RuneCountInString(summary)
	if count > 4000 {
		s := []rune(summary)
		summary = string(s[0:4000])
	}
	i.ISummary = &ISummary{
		Text: summary,
	}
}

// AddDuration adds the duration to the iTunes duration field.
func (i *Item) AddDuration(durationInSeconds int64) {
	if durationInSeconds <= 0 {
		return
	}
	i.IDuration = parseDuration(durationInSeconds)
}

// String returns the MIME type encoding of the specified EnclosureType.
func (et EnclosureType) String() string {
	// https://help.apple.com/itc/podcasts_connect/#/itcb54353390
	switch et {
	case M4A:
		return "audio/x-m4a"
	case M4V:
		return "video/x-m4v"
	case MP4:
		return "video/mp4"
	case MP3:
		return "audio/mpeg"
	case MOV:
		return "video/quicktime"
	case PDF:
		return "application/pdf"
	case EPUB:
		return "document/x-epub"
	}
	return enclosureDefault
}
