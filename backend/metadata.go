package backend

import (
	"context"
	"strings"

	"cloud.google.com/go/datastore"
	ds "github.com/txsvc/platform/v2/pkg/datastore"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/errordef"
	"github.com/podops/podops/internal/metadata"
)

const (
	// DatastoreMetadata collection METADATA
	datastoreMetadata = "METADATA"
)

// GetMetadata retrieves the metadata for a resource
func GetMetadata(ctx context.Context, guid string) (*metadata.Metadata, error) {
	var m metadata.Metadata

	if err := ds.DataStore().Get(ctx, metadataKey(guid), &m); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil // not found is not an error
		}
		return nil, err
	}
	return &m, nil
}

// GetMetadataForResource retrieves the metadata associated with a resource, if the resource is of type "episode", nil otherwise.
func GetMetadataForResource(ctx context.Context, guid string) (*metadata.Metadata, error) {
	r, err := GetResource(ctx, guid)
	if err != nil {
		return nil, err
	}
	if r.Kind == podops.ResourceAsset {
		return nil, nil
	}

	metaGUID := ""
	if r.Kind == podops.ResourceShow {
		metaGUID = strings.Split(metadata.LocalNamePart(r.ImageURI), ".")[0]
	} else if r.Kind == podops.ResourceEpisode {
		metaGUID = strings.Split(metadata.LocalNamePart(r.EnclosureURI), ".")[0]
	} else {
		return nil, errordef.ErrNoSuchResource
	}

	return GetMetadata(ctx, metaGUID)
}

// UpdateMetadata does what the name suggests
func UpdateMetadata(ctx context.Context, m *metadata.Metadata) error {
	if _, err := ds.DataStore().Put(ctx, metadataKey(m.GUID), m); err != nil {
		return err
	}
	return nil
}

// DeleteMetadata deletes a resource and it's backing .yaml file
func DeleteMetadata(ctx context.Context, guid string) error {
	m, err := GetMetadata(ctx, guid)
	if err != nil {
		return err
	}
	if m == nil { // not found
		return errordef.ErrNoSuchResource
	}

	if err := ds.DataStore().Delete(ctx, metadataKey(m.GUID)); err != nil {
		return err
	}
	return nil
}

func metadataKey(guid string) *datastore.Key {
	return datastore.NameKey(datastoreMetadata, guid, nil)
}
