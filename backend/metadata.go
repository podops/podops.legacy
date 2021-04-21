package backend

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/fupas/platform/pkg/platform"

	"github.com/podops/podops/internal/errordef"
	"github.com/podops/podops/internal/metadata"
)

const (
	// DatastoreMetadata collection METADATA
	DatastoreMetadata = "METADATA"
)

// GetResourceMetadata retrieves the metadata for a resource
func GetResourceMetadata(ctx context.Context, guid string) (*metadata.Metadata, error) {
	var m metadata.Metadata

	if err := platform.DataStore().Get(ctx, metadataKey(guid), &m); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil // not found is not an error
		}
		return nil, err
	}
	return &m, nil
}

// UpdateResourceMetadata does what the name suggests
func UpdateResourceMetadata(ctx context.Context, m *metadata.Metadata) error {
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

func metadataKey(guid string) *datastore.Key {
	return datastore.NameKey(DatastoreMetadata, guid, nil)
}
