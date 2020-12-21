package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
)

var parseDuration = func(duration int64) string {
	h := duration / 3600
	duration = duration % 3600

	m := duration / 60
	duration = duration % 60

	s := duration

	// HH:MM:SS
	if h > 9 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}

	// H:MM:SS
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}

	// MM:SS
	if m > 9 {
		return fmt.Sprintf("%02d:%02d", m, s)
	}

	// M:SS
	return fmt.Sprintf("%d:%02d", m, s)
}

var parseDateRFC1123Z = func(t *time.Time) string {
	if t != nil && !t.IsZero() {
		return t.Format(time.RFC1123Z)
	}
	return time.Now().UTC().Format(time.RFC1123Z)
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

var encoder = func(w io.Writer, o interface{}) error {
	e := xml.NewEncoder(w)
	e.Indent("", "  ")
	if err := e.Encode(o); err != nil {
		return errors.Wrap(err, "channel.encoder: e.Encode returned error")
	}
	return nil
}

var parseAuthorNameEmail = func(a *Author) string {
	var author string
	if a != nil {
		author = a.Email
		if len(a.Name) > 0 {
			author = fmt.Sprintf("%s (%s)", a.Email, a.Name)
		}
	}
	return author
}
