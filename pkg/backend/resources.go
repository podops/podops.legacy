package backend

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/platform"
	a "github.com/podops/podops/apiv1"
	"gopkg.in/yaml.v2"
)

const (
	// DatastoreResources collection RESOURCE
	DatastoreResources = "RESOURCES"
)

var (
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
		return "", fmt.Errorf("invalid resource '%s'", kind)
	}
	return k, nil
}

// GetResource retrieves a resource
func GetResource(ctx context.Context, guid string) (*a.Resource, error) {
	var r a.Resource

	if err := platform.DataStore().Get(ctx, resourceKey(guid), &r); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil // not found is not an error
		}
		return nil, err
	}
	return &r, nil
}

// FindResource looks for a resource 'name' in the context of production 'parent'
func FindResource(ctx context.Context, production, name string) (*a.Resource, error) {
	var r []*a.Resource

	if _, err := platform.DataStore().GetAll(ctx, datastore.NewQuery(DatastoreResources).Filter("ParentGUID =", production).Filter("Name =", name), &r); err != nil {
		return nil, err
	}
	if r == nil {
		return nil, nil
	}
	if len(r) > 1 {
		return nil, fmt.Errorf("inconsistent inventory: expected 1, found %d resources", len(r))
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
			return fmt.Errorf("can not update resource: expected '%s', received '%s'", r.Kind, _kind)
		}
		r.Name = name
		r.ParentGUID = production
		r.Location = location
		r.Updated = util.Timestamp()

		return updateResource(ctx, r)
	}

	// create a new inventory entry
	now := util.Timestamp()
	rsrc := a.Resource{
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

// UpdateAssetResource updates the resource inventory
func UpdateAssetResource(ctx context.Context, name, guid, kind, production, location, contentType, original, etag string, size, duration int64) error {
	r, _ := GetResource(ctx, guid)

	_kind, err := NormalizeKind(kind)
	if err != nil {
		return err
	}

	if r != nil {
		// resource already exists, just update the inventory
		if r.Kind != _kind {
			return fmt.Errorf("can not modify resource: expected '%s', received '%s'", r.Kind, _kind)
		}
		r.Name = name
		r.ParentGUID = production
		r.Location = location
		r.Extra1 = original
		r.Extra2 = etag
		r.ContentType = contentType
		r.Size = size
		r.Duration = duration
		r.Updated = util.Timestamp()

		return updateResource(ctx, r)
	}

	// create a new inventory entry
	now := util.Timestamp()
	rsrc := a.Resource{
		Name:        name,
		GUID:        guid,
		Kind:        _kind,
		ParentGUID:  production,
		Location:    location,
		Extra1:      original,
		Extra2:      etag,
		ContentType: contentType,
		Size:        size,
		Duration:    duration,
		Created:     now,
		Updated:     now,
	}
	return updateResource(ctx, &rsrc)
}

// DeleteResource deletes a resource and it's backing .yaml file
func DeleteResource(ctx context.Context, guid string) error {
	r, err := GetResource(ctx, guid)
	if err != nil {
		return err
	}
	if r == nil { // not found
		return a.ErrNoSuchResource
	}

	if err := platform.DataStore().Delete(ctx, resourceKey(r.GUID)); err != nil {
		return err // FIXME put r back if this fails?
	}

	// validate the production after deleting a resource
	prod := ""
	if r.Kind == a.ResourceShow {
		prod = r.GUID
	} else {
		prod = r.ParentGUID
	}
	if err = ValidateProduction(ctx, prod); err != nil {
		p, err := GetProduction(ctx, prod)
		if err != nil {
			return err
		}
		p.BuildDate = 0 // FIXME BuildDate is the only flag we currently have to mark a production as VALID
		UpdateProduction(ctx, p)
	}

	if r.Kind == a.ResourceAsset {
		return RemoveAsset(ctx, r.Location)
	}
	return RemoveResource(ctx, r.Location)
}

// ListResources returns all resources of type kind belonging to parentID
func ListResources(ctx context.Context, production, kind string) ([]*a.Resource, error) {
	var r []*a.Resource

	_kind, err := NormalizeKind(kind)
	if err != nil {
		return nil, err
	}

	if _kind == a.ResourceALL {
		if _, err := platform.DataStore().GetAll(ctx, datastore.NewQuery(DatastoreResources).Filter("ParentGUID =", production).Order("-Created"), &r); err != nil {
			return nil, err
		}
		// as we do not get SHOW with the query, we add it now
		show, err := GetResource(ctx, production)
		if err == nil && show != nil { // SHOW could not be there, no worries ...
			r = append(r, show)
		}
	} else if _kind == a.ResourceShow {
		// there should only be ONE
		show, err := GetResource(ctx, production)
		if err == nil && show != nil { // SHOW could not be there, no worries ...
			r = append(r, show)
		}
	} else {
		if _, err := platform.DataStore().GetAll(ctx, datastore.NewQuery(DatastoreResources).Filter("ParentGUID =", production).Filter("Kind =", _kind).Order("-Created"), &r); err != nil {
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

	if r.Kind == a.ResourceAsset {
		asset := a.Asset{
			URI:   fmt.Sprintf("%s/%s", a.DefaultCDNEndpoint, r.Location),
			Title: r.Extra1,
			Type:  r.ContentType,
			Size:  int(r.Size),
			Rel:   a.ResourceTypeLocal,
		}
		return &asset, nil
	}

	rsrc, _, _, err := ReadResource(ctx, r.Location)
	if err != nil {
		return nil, err
	}

	return rsrc, nil
}

// WriteResourceContent creates a resource .yaml file. An existing resource will be overwritten if force==true
func WriteResourceContent(ctx context.Context, path string, create, force bool, rsrc interface{}) error {

	exists := true

	bkt := platform.Storage().Bucket(a.BucketProduction)
	obj := bkt.Object(path)

	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		exists = false
	}

	// some logic mangling here ...
	if create && exists && !force { // create on an existing resource
		return fmt.Errorf("'%s' already exists", path)
	}
	if !exists && !create && !force { // update on a missing resource
		return fmt.Errorf("'%s' does not exists", path)
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

// ReadResource reads a resource from Cloud Storage
func ReadResource(ctx context.Context, path string) (interface{}, string, string, error) {

	bkt := platform.Storage().Bucket(a.BucketProduction)
	reader, err := bkt.Object(path).NewReader(ctx)
	if err != nil {
		return nil, "", "", err
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, "", "", err
	}

	return a.LoadResource(data)
}

// RemoveResource removes a resource from Cloud Storage
func RemoveResource(ctx context.Context, path string) error {
	bkt := platform.Storage().Bucket(a.BucketProduction)

	obj := bkt.Object(path)
	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return a.ErrNoSuchResource
	}

	return bkt.Object(path).Delete(ctx)
}

// RemoveAsset removes a asset from Cloud Storage
func RemoveAsset(ctx context.Context, path string) error {
	bkt := platform.Storage().Bucket(a.BucketCDN)

	obj := bkt.Object(path)
	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return a.ErrNoSuchAsset
	}

	return bkt.Object(path).Delete(ctx)
}

// UpdateShow is a helper function to update a show resource
func UpdateShow(ctx context.Context, location string, show *a.Show) error {
	r, _ := GetResource(ctx, show.GUID())

	if r != nil {
		// resource already exists, just update the inventory
		if r.Kind != show.Kind {
			return fmt.Errorf("can not modify resource: expected '%s', received '%s'", r.Kind, show.Kind)
		}
		r.Name = show.Metadata.Name
		r.ParentGUID = show.Metadata.Labels[a.LabelParentGUID]
		r.Location = location
		r.Title = show.Description.Title
		r.Summary = show.Description.Summary
		// FIXME: r.Image = show.Image.ResolveURI(a.DefaultCDNEndpoint, show.GUID())
		r.Image = show.Image.ResolveURI(a.StorageEndpoint, show.GUID())
		r.Updated = util.Timestamp()

		return updateResource(ctx, r)
	}

	// create a new inventory entry
	now := util.Timestamp()
	rsrc := a.Resource{
		Name:       show.Metadata.Name,
		GUID:       show.GUID(),
		Kind:       a.ResourceShow,
		ParentGUID: show.Metadata.Labels[a.LabelParentGUID],
		Location:   location,
		Title:      show.Description.Title,
		Summary:    show.Description.Summary,
		// FIXME: Image:      show.Image.ResolveURI(a.DefaultCDNEndpoint, show.GUID()),
		Image:   show.Image.ResolveURI(a.StorageEndpoint, show.GUID()),
		Created: now,
		Updated: now,
	}
	return updateResource(ctx, &rsrc)
}

// UpdateEpisode is a helper function to update a episode resource
func UpdateEpisode(ctx context.Context, location string, episode *a.Episode) error {
	// check if resource with same name already exists for the parent production
	rn, err := FindResource(ctx, episode.ParentGUID(), episode.Metadata.Name)
	if err != nil {
		return err
	}
	r, err := GetResource(ctx, episode.GUID())
	if err != nil {
		return err
	}

	if rn != nil && r != nil {
		if rn.GUID != r.GUID {
			return fmt.Errorf("can not update resource: '%s/%s' already exists", episode.ParentGUID(), episode.Metadata.Name)
		}
	}

	if r != nil {
		// resource already exists, just update the inventory
		if r.Kind != episode.Kind {
			return fmt.Errorf("can not modify resource: expected '%s', received '%s'", r.Kind, episode.Kind)
		}
		index, _ := strconv.ParseInt(episode.Metadata.Labels[a.LabelEpisode], 10, 64)

		r.Name = episode.Metadata.Name
		r.ParentGUID = episode.Metadata.Labels[a.LabelParentGUID]
		r.Location = location
		r.Title = episode.Description.Title
		r.Summary = episode.Description.Summary
		r.Published = episode.PublishDateTimestamp()
		r.Index = int(index) // episode number
		r.Image = episode.Image.ResolveURI(a.StorageEndpoint, episode.ParentGUID())
		r.Extra1 = episode.Enclosure.ResolveURI(a.DefaultCDNEndpoint+"/c", episode.ParentGUID())
		r.Size = int64(episode.Enclosure.Size)
		r.Duration = int64(episode.Description.Duration)
		r.Updated = util.Timestamp()

		return updateResource(ctx, r)
	}

	// create a new inventory entry
	now := util.Timestamp()
	index, _ := strconv.ParseInt(episode.Metadata.Labels[a.LabelEpisode], 10, 64)

	rsrc := a.Resource{
		Name:       episode.Metadata.Name,
		GUID:       episode.GUID(),
		Kind:       a.ResourceEpisode,
		ParentGUID: episode.Metadata.Labels[a.LabelParentGUID],
		Location:   location,
		Title:      episode.Description.Title,
		Summary:    episode.Description.Summary,
		Published:  episode.PublishDateTimestamp(),
		Index:      int(index), // episode number
		Image:      episode.Image.ResolveURI(a.StorageEndpoint, episode.ParentGUID()),
		Extra1:     episode.Enclosure.ResolveURI(a.DefaultCDNEndpoint+"/c", episode.ParentGUID()),
		Size:       int64(episode.Enclosure.Size),
		Duration:   int64(episode.Description.Duration),
		Created:    now,
		Updated:    now,
	}
	return updateResource(ctx, &rsrc)
}

// updateResource does what the name suggests
func updateResource(ctx context.Context, r *a.Resource) error {
	if _, err := platform.DataStore().Put(ctx, resourceKey(r.GUID), r); err != nil {
		return err
	}
	return nil
}

func resourceKey(guid string) *datastore.Key {
	return datastore.NameKey(DatastoreResources, guid, nil)
}
