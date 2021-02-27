package rss

import (
	"encoding/xml"
	"io"
	"time"
)

type (
	// Channel represents a RSS feed with iTunes/podcast specific extensions
	Channel struct {
		XMLName        xml.Name `xml:"channel"`
		Title          string   `xml:"title"`
		Link           string   `xml:"link"`
		Description    string   `xml:"description"`
		Category       string   `xml:"category,omitempty"`
		Cloud          string   `xml:"cloud,omitempty"`
		Copyright      string   `xml:"copyright,omitempty"`
		Docs           string   `xml:"docs,omitempty"`
		Generator      string   `xml:"generator,omitempty"`
		Language       string   `xml:"language,omitempty"`
		LastBuildDate  string   `xml:"lastBuildDate,omitempty"`
		ManagingEditor string   `xml:"managingEditor,omitempty"`
		PubDate        string   `xml:"pubDate,omitempty"`
		Rating         string   `xml:"rating,omitempty"`
		SkipHours      string   `xml:"skipHours,omitempty"`
		SkipDays       string   `xml:"skipDays,omitempty"`
		TTL            int      `xml:"ttl,omitempty"`
		WebMaster      string   `xml:"webMaster,omitempty"`
		Image          *Image
		TextInput      *TextInput
		AtomLink       *AtomLink

		// https://help.apple.com/itc/podcasts_connect/#/itcb54353390
		IAuthor     string `xml:"itunes:author,omitempty"`
		ISubtitle   string `xml:"itunes:subtitle,omitempty"`
		ISummary    *ISummary
		IBlock      string `xml:"itunes:block,omitempty"`
		IImage      *IImage
		IDuration   string  `xml:"itunes:duration,omitempty"`
		IExplicit   string  `xml:"itunes:explicit,omitempty"`
		IComplete   string  `xml:"itunes:complete,omitempty"`
		INewFeedURL string  `xml:"itunes:new-feed-url,omitempty"`
		IOwner      *Author // Author is formatted for itunes as-is
		ICategories []*ICategory
		// ADDED
		ITitle string `xml:"itunes:title,omitempty"`
		IType  string `xml:"itunes:type,omitempty"`

		Items []*Item

		encode func(w io.Writer, o interface{}) error
	}

	channelWrapper struct {
		XMLName  xml.Name `xml:"rss"`
		Version  string   `xml:"version,attr"`
		ATOMNS   string   `xml:"xmlns:atom,attr,omitempty"`
		ITUNESNS string   `xml:"xmlns:itunes,attr"`
		Channel  *Channel
	}

	// Item represents a single entry in a podcast.
	//
	// Article minimal requirements are:
	// - Title
	// - Description
	// - Link
	//
	// Audio minimal requirements are:
	// - Title
	// - Description
	// - Enclosure (HREF, Type and Length all required)
	//
	// Recommendations:
	// - Setting the minimal fields sets most of other fields, including iTunes.
	// - Use the Published time.Time setting instead of PubDate.
	// - Always set an Enclosure.Length, to be nice to your downloaders.
	// - Use Enclosure.Type instead of setting TypeFormatted for valid extensions.
	//
	Item struct {
		XMLName          xml.Name   `xml:"item"`
		GUID             string     `xml:"guid"`
		Title            string     `xml:"title"`
		Link             string     `xml:"link"`
		Description      string     `xml:"description"`
		Author           *Author    `xml:"-"`
		AuthorFormatted  string     `xml:"author,omitempty"`
		Category         string     `xml:"category,omitempty"`
		Comments         string     `xml:"comments,omitempty"`
		Source           string     `xml:"source,omitempty"`
		PubDate          *time.Time `xml:"-"`
		PubDateFormatted string     `xml:"pubDate,omitempty"`
		Enclosure        *Enclosure

		// https://help.apple.com/itc/podcasts_connect/#/itcb54353390
		IAuthor   string `xml:"itunes:author,omitempty"`
		ISubtitle string `xml:"itunes:subtitle,omitempty"`
		ISummary  *ISummary
		IImage    *IImage
		IDuration string `xml:"itunes:duration,omitempty"`
		IExplicit string `xml:"itunes:explicit,omitempty"`
		// ADDED
		ISeason      string `xml:"itunes:season,omitempty"`
		IEpisode     string `xml:"itunes:episode,omitempty"`
		IEpisodeType string `xml:"itunes:episodeType,omitempty"`
		IBlock       string `xml:"itunes:block,omitempty"`

		// REMOVE IIsClosedCaptioned string `xml:"itunes:isClosedCaptioned,omitempty"`
		// REMOVE IOrder string `xml:"itunes:order,omitempty"`
	}

	// AtomLink represents the Atom reference link.
	AtomLink struct {
		XMLName xml.Name `xml:"atom:link"`
		HREF    string   `xml:"href,attr"`
		Rel     string   `xml:"rel,attr"`
		Type    string   `xml:"type,attr"`
	}

	// Image represents an image.
	//
	// Podcast feeds contain artwork that is a minimum size of 1400 x 1400 pixels and a maximum size of 3000 x 3000 pixels,
	// 72 dpi, in JPEG or PNG format with appropriate file extensions (.jpg, .png), and in the RGB colorspace. To optimize
	// images for mobile devices, Apple recommends compressing your image files.
	Image struct {
		XMLName     xml.Name `xml:"image"`
		URL         string   `xml:"url"`
		Title       string   `xml:"title"`
		Link        string   `xml:"link"`
		Description string   `xml:"description,omitempty"`
		Width       int      `xml:"width,omitempty"`
		Height      int      `xml:"height,omitempty"`
	}

	// Author represents a named author and email.
	//
	// For iTunes compliance, both Name and Email are required.
	Author struct {
		XMLName xml.Name `xml:"itunes:owner"`
		Name    string   `xml:"itunes:name"`
		Email   string   `xml:"itunes:email"`
	}

	// TextInput represents text inputs.
	TextInput struct {
		XMLName     xml.Name `xml:"textInput"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
		Name        string   `xml:"name"`
		Link        string   `xml:"link"`
	}

	// ICategory is a 2-tier classification system for iTunes.
	ICategory struct {
		XMLName     xml.Name `xml:"itunes:category"`
		Text        string   `xml:"text,attr"`
		ICategories []*ICategory
	}

	// IImage represents an iTunes image.
	//
	// Podcast feeds contain artwork that is a minimum size of
	// 1400 x 1400 pixels and a maximum size of 3000 x 3000 pixels,
	// 72 dpi, in JPEG or PNG format with appropriate file
	// extensions (.jpg, .png), and in the RGB colorspace. To optimize
	// images for mobile devices, Apple recommends compressing your
	// image files.
	IImage struct {
		XMLName xml.Name `xml:"itunes:image"`
		HREF    string   `xml:"href,attr"`
	}

	// ISummary is a 4000 character rich-text field for the itunes:summary tag.
	//
	// This is rendered as CDATA which allows for HTML tags such as `<a href="">`.
	ISummary struct {
		XMLName xml.Name `xml:"itunes:summary"`
		Text    string   `xml:",cdata"`
	}

	// EnclosureType specifies the type of the enclosure.
	EnclosureType int

	// Enclosure represents a download enclosure.
	Enclosure struct {
		XMLName xml.Name `xml:"enclosure"`

		// URL is the downloadable url for the content. (Required)
		URL string `xml:"url,attr"`

		// Length is the size in Bytes of the download. (Required)
		Length int64 `xml:"-"`
		// LengthFormatted is the size in Bytes of the download. (Required)
		//
		// This field gets overwritten with the API when setting Length.
		LengthFormatted string `xml:"length,attr"`

		// Type is MIME type encoding of the download. (Required)
		Type EnclosureType `xml:"-"`
		// TypeFormatted is MIME type encoding of the download. (Required)
		//
		// This field gets overwritten with the API when setting Type.
		TypeFormatted string `xml:"type,attr"`
	}
)
