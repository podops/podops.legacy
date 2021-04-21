package metadata

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tcolgate/mp3"
)

const (
	defaultContentType = "application/octet-stream"
)

type (
	// Metadata keeps basic metadata of a cdn resource
	Metadata struct {
		Name        string `json:"name"`
		GUID        string `json:"guid"`
		Size        int64  `json:"size"`
		Duration    int64  `json:"duration"`
		ContentType string `json:"content_type"`
		Etag        string `json:"etag"`
		Timestamp   int64  `json:"timestamp"`
	}
)

// ExtractMetadataFromResponse extracts the metadata from http.Response
func ExtractMetadataFromResponse(resp *http.Response) *Metadata {
	if resp == nil {
		return nil
	}

	meta := Metadata{
		ContentType: resp.Header.Get("content-type"),
		Etag:        resp.Header.Get("etag"),
	}
	l, err := strconv.ParseInt(resp.Header.Get("content-length"), 10, 64)
	if err == nil {
		meta.Size = l
	}
	// expects 'Wed, 30 Dec 2020 14:14:26 GM'
	t, err := time.Parse(time.RFC1123, resp.Header.Get("date"))
	if err == nil {
		meta.Timestamp = t.Unix()
	}
	return &meta
}

func ExtractMetadataFromFile(path string) (*Metadata, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// the basics
	meta := Metadata{
		Name:        fi.Name(),
		Size:        fi.Size(),
		ContentType: defaultContentType,
		Timestamp:   fi.ModTime().Unix(),
	}

	// calculate our etag
	meta.Etag = meta.ETAG()

	// try to detect the media type
	// thanks to https://gist.github.com/rayrutjes/db9b9ea8e02255d62ce2
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}
	meta.ContentType = http.DetectContentType(buffer)
	// reset the read pointer
	file.Seek(0, 0)

	// in case it is a .mp3, calculate the play time.
	// thanks to https://stackoverflow.com/questions/60281655/how-to-find-the-length-of-mp3-file-in-golang
	if meta.IsMP3() {
		d := mp3.NewDecoder(file)

		var f mp3.Frame
		skipped := 0
		t := 0.0

		for {
			if err := d.Decode(&f, &skipped); err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
			t = t + f.Duration().Seconds()
		}
		meta.Duration = int64(t) // duration in seconds
	}
	return &meta, nil
}

func (m *Metadata) ETAG() string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s%d%d", m.Name, m.Size, m.Timestamp)))
	return hex.EncodeToString(hash[:])
}

func (m *Metadata) IsMP3() bool {
	return m.ContentType == "audio/mpeg"
}

// CalculateLength returns the play duration of a media file like a .mp3
func CalculateLength(path string) (int64, error) {
	m, err := ExtractMetadataFromFile(path)
	if err != nil {
		return 0, err
	}
	return m.Duration, nil
}
