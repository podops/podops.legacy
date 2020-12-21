package rss

import (
	"encoding/xml"
)

// iTunes Specifications: https://help.apple.com/itc/podcasts_connect/#/itcb54353390

// EnclosureType specifies the type of the enclosure.
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

type (

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
