package backend

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/fupas/platform/pkg/platform"
	"github.com/podops/podops"
	"github.com/podops/podops/internal/errordef"
)

const (
	// DatastoreMetadata collection METADATA
	DatastoreMetadata = "METADATA"
)

// GetResourceMetadata retrieves the metadata for a resource
func GetResourceMetadata(ctx context.Context, guid string) (*podops.ResourceMetadata, error) {
	var m podops.ResourceMetadata

	if err := platform.DataStore().Get(ctx, metadataKey(guid), &m); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil // not found is not an error
		}
		return nil, err
	}
	return &m, nil
}

// UpdateResourceMetadata does what the name suggests
func UpdateResourceMetadata(ctx context.Context, m *podops.ResourceMetadata) error {
	if _, err := platform.DataStore().Put(ctx, metadataKey(m.GUID), m); err != nil {
		return err
	}
	return nil
}

// DeleteResource deletes a resource and it's backing .yaml file
func DeleteResourceMetadata(ctx context.Context, guid string) error {
	m, err := GetResourceMetadata(ctx, guid)
	if err != nil {
		return err
	}
	if m == nil { // not found
		return errordef.ErrNoSuchResource
	}

	if err := platform.DataStore().Delete(ctx, metadataKey(m.GUID)); err != nil {
		return err
	}
	return nil
}

// ExtractMetadataFromResponse extracts the metadata from http.Response
func ExtractMetadataFromResponse(resp *http.Response) *podops.ResourceMetadata {
	meta := podops.ResourceMetadata{
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

// CalculateLength returns the play duration of a media file like a .mp3
func CalculateLength(contentType, path string) int64 {
	return 0 // FIXME to be implemented
}

func metadataKey(guid string) *datastore.Key {
	return datastore.NameKey(DatastoreMetadata, guid, nil)
}
