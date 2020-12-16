package feed

import "encoding/xml"

// AtomLink represents the Atom reference link.
type AtomLink struct {
	XMLName xml.Name `xml:"atom:link"`
	HREF    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
}

// Image represents an image.
//
// Podcast feeds contain artwork that is a minimum size of
// 1400 x 1400 pixels and a maximum size of 3000 x 3000 pixels,
// 72 dpi, in JPEG or PNG format with appropriate file
// extensions (.jpg, .png), and in the RGB colorspace. To optimize
// images for mobile devices, Apple recommends compressing your
// image files.
type Image struct {
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
type Author struct {
	XMLName xml.Name `xml:"itunes:owner"`
	Name    string   `xml:"itunes:name"`
	Email   string   `xml:"itunes:email"`
}

// TextInput represents text inputs.
type TextInput struct {
	XMLName     xml.Name `xml:"textInput"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Name        string   `xml:"name"`
	Link        string   `xml:"link"`
}
