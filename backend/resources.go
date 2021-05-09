package backend

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"

	"gopkg.in/yaml.v2"

	ds "github.com/txsvc/platform/v2/pkg/datastore"
	"github.com/txsvc/platform/v2/pkg/timestamp"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/errordef"
	"github.com/podops/podops/internal/loader"
	"github.com/podops/podops/internal/messagedef"
)

const (
	// DatastoreResources collection RESOURCE
	datastoreResources = "RESOURCES"
)

var (
	// full canonical route
	importTaskEndpoint string = podops.DefaultCDNEndpoint + "/_w/import"
	syncTaskEndpoint   string = podops.DefaultCDNEndpoint + "/_w/sync"
	// mapping of resource names and aliases
	resourceMap map[string]string
)

func init() {
	resourceMap = make(map[string]string)
	resourceMap["show"] = "show"
	resourceMap["shows"] = "show"
	resourceMap["episode"] = "episode"
	resourceMap["episodes"] = "episode"
	resourceMap["asset"] = "asset"
	resourceMap["assets"] = "asset"
	resourceMap["all"] = "all"
}

func NormalizeKind(kind string) (string, error) {
	k := resourceMap[strings.ToLower(kind)]
	if k == "" {
		return "", fmt.Errorf(messagedef.MsgResourceIsInvalid, kind)
	}
	return k, nil
}

// GetResource retrieves a resource
func GetResource(ctx context.Context, guid string) (*podops.Resource, error) {
	var r podops.Resource

	if err := ds.DataStore().Get(ctx, resourceKey(guid), &r); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil // not found is not an error
		}
		return nil, err
	}
	return &r, nil
}

// FindResource looks for a resource 'name' in the context of production 'production'
func FindResource(ctx context.Context, production, name string) (*podops.Resource, error) {
	var r []*podops.Resource

	if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(datastoreResources).Filter("ParentGUID =", production).Filter("Name =", name), &r); err != nil {
		return nil, err
	}
	if r == nil {
		return nil, nil
	}
	if len(r) > 1 {
		return nil, fmt.Errorf(messagedef.MsgResourceInconsistentInventory, 1, len(r), production, name)
	}

	return r[0], nil
}

// UpdateResource updates the resource inventory
func UpdateResource(ctx context.Context, name, guid, kind, production, location string) error {
	r, _ := GetResource(ctx, guid)

	_kind, err := NormalizeKind(kind)
	if err != nil {
		return err
	}

	if r != nil {
		// resource already exists, just update the inventory
		if r.Kind != _kind {
			return fmt.Errorf(messagedef.MsgResourceKindMismatch, r.Kind, _kind)
		}
		r.Name = name
		r.ParentGUID = production
		r.Location = location
		r.Updated = timestamp.Now()

		return updateResource(ctx, r)
	}

	// create a new inventory entry
	now := timestamp.Now()
	rsrc := podops.Resource{
		Name:       name,
		GUID:       guid,
		Kind:       _kind,
		ParentGUID: production,
		Location:   location,
		Created:    now,
		Updated:    now,
	}
	return updateResource(ctx, &rsrc)
}

// DeleteResource deletes a resource and it's backing .yaml file
func DeleteResource(ctx context.Context, prod, kind, guid string) error {
	r, err := GetResource(ctx, guid)
	if err != nil {
		return err
	}
	if r == nil { // not found
		return errordef.ErrNoSuchResource
	}

	if err := ds.DataStore().Delete(ctx, resourceKey(r.GUID)); err != nil {
		return err
	}

	// validate the production after deleting a resource
	if err = ValidateProduction(ctx, prod); err != nil {
		p, err := GetProduction(ctx, prod)
		if err != nil {
			return err
		}
		p.BuildDate = 0
		p.Published = false
		p.LatestPublishDate = 0
		UpdateProduction(ctx, p)
	}

	if r.Kind == podops.ResourceAsset {
		if err := DeleteMetadata(ctx, guid); err != nil {
			return err
		}
		return RemoveAsset(ctx, prod, r.Location)
	}
	return RemoveResourceContent(ctx, r.Location)
}

// ListResources returns all resources of type kind belonging to parentID
func ListResources(ctx context.Context, production, kind string) ([]*podops.Resource, error) {
	var r []*podops.Resource

	_kind, err := NormalizeKind(kind)
	if err != nil {
		return nil, err
	}

	if _kind == podops.ResourceALL {
		if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(datastoreResources).Filter("ParentGUID =", production).Order("-Created"), &r); err != nil {
			return nil, err
		}
	} else if _kind == podops.ResourceShow {
		// there should only be ONE
		show, err := GetResource(ctx, production)
		if err == nil && show != nil { // SHOW could not be there, no worries ...
			r = append(r, show)
		}
	} else {
		if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(datastoreResources).Filter("ParentGUID =", production).Filter("Kind =", _kind).Order("-Created"), &r); err != nil {
			return nil, err
		}
	}

	if len(r) == 0 {
		return nil, nil
	}
	return r, nil
}

// GetResourceContent retrieves a resource file
func GetResourceContent(ctx context.Context, guid string) (interface{}, error) {
	r, err := GetResource(ctx, guid)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, nil // not found => not an eror
	}

	if r.Kind == podops.ResourceAsset {
		meta, err := GetMetadata(ctx, guid)
		if err != nil {
			return nil, err
		}

		asset := podops.Asset{
			URI:   r.GetPublicLocation(),
			Title: r.Name,
			Type:  meta.ContentType,
			Size:  int(meta.Size), // GITHUB_ISSUE #10
			Rel:   podops.ResourceTypeLocal,
		}
		return &asset, nil
	}

	rsrc, _, _, err := ReadResourceContent(ctx, r.Location)
	if err != nil {
		return nil, err
	}

	// metadata mix-in
	if r.Kind == podops.ResourceShow {
		show := rsrc.(*podops.Show)
		show.Image.URI = r.ImageURI
		return show, nil

	} else if r.Kind == podops.ResourceEpisode {
		meta, err := GetMetadataForResource(ctx, guid)
		if err != nil {
			return nil, err
		}
		episode := rsrc.(*podops.Episode)
		episode.Enclosure.URI = r.EnclosureURI
		episode.Enclosure.Size = int(meta.Size) // GITHUB_ISSUE #10
		episode.Description.Duration = int(meta.Duration)

		return episode, nil
	}

	return nil, errordef.ErrNoSuchResource
}

// WriteResourceContent creates a resource .yaml file. An existing resource will be overwritten if force==true
func WriteResourceContent(ctx context.Context, path string, create, force bool, rsrc interface{}) error {

	exists := true

	bkt := ds.Storage().Bucket(podops.BucketProduction)
	obj := bkt.Object(path)

	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		exists = false
	}

	// some logic mangling here ...
	if create && exists && !force { // create on an existing resource
		return fmt.Errorf(messagedef.MsgResourceAlreadyExists, path)
	}
	if !exists && !create && !force { // update on a missing resource
		return fmt.Errorf(messagedef.MsgResourceNotFound, path)
	}

	data, err := yaml.Marshal(rsrc)
	if err != nil {
		return err
	}

	writer := obj.NewWriter(ctx)
	defer writer.Close()
	if _, err := writer.Write(data); err != nil {
		return err
	}

	return nil
}

// ReadResourceContent reads a resource from the Cloud Storage
func ReadResourceContent(ctx context.Context, path string) (interface{}, string, string, error) {

	bkt := ds.Storage().Bucket(podops.BucketProduction)
	reader, err := bkt.Object(path).NewReader(ctx)
	if err != nil {
		return nil, "", "", err
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, "", "", err
	}

	return loader.UnmarshalResource(data)
}

// RemoveResourceContent removes a resource from Cloud Storage
func RemoveResourceContent(ctx context.Context, location string) error {
	bkt := ds.Storage().Bucket(podops.BucketProduction)

	obj := bkt.Object(location)
	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return errordef.ErrNoSuchResource
	}

	return bkt.Object(location).Delete(ctx)
}

// updateResource does what the name suggests
func updateResource(ctx context.Context, r *podops.Resource) error {
	if _, err := ds.DataStore().Put(ctx, resourceKey(r.GUID), r); err != nil {
		return err
	}
	return nil
}

func resourceKey(guid string) *datastore.Key {
	return datastore.NameKey(datastoreResources, guid, nil)
}
