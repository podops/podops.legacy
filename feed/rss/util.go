package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"
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

var encoder = func(w io.Writer, o interface{}) error {
	e := xml.NewEncoder(w)
	e.Indent("", "  ")
	if err := e.Encode(o); err != nil {
		return fmt.Errorf("channel.encoder: error %v", err)
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
