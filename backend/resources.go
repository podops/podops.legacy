package backend

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"gopkg.in/yaml.v2"

	"github.com/fupas/commons/pkg/env"
	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/platform"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/errordef"
	"github.com/podops/podops/internal/loader"
	"github.com/podops/podops/internal/metadata"
	p "github.com/podops/podops/internal/platform"
)

const (
	// DatastoreResources collection RESOURCE
	DatastoreResources = "RESOURCES"
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
		return "", fmt.Errorf("invalid resource '%s'", kind)
	}
	return k, nil
}

// GetResource retrieves a resource
func GetResource(ctx context.Context, guid string) (*podops.Resource, error) {
	var r podops.Resource

	if err := platform.DataStore().Get(ctx, resourceKey(guid), &r); err != nil {
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

// UpdateAsset updates the resource inventory
func UpdateAsset(ctx context.Context, meta *metadata.Metadata, production, location, rel string) error {
	r, _ := GetResource(ctx, meta.GUID)

	if r != nil {
		// resource already exists, just update the inventory
		r.Name = meta.Name
		r.ParentGUID = production
		r.Location = location
		r.Updated = util.Timestamp()

		if meta.IsImage() {
			r.ImageURI = meta.Origin
			r.ImageRel = rel
		} else {
			r.EnclosureURI = meta.Origin
			r.EnclosureRel = rel
		}

		if err := UpdateMetadata(ctx, meta); err != nil {
			return err
		}
		return updateResource(ctx, r)
	}

	// create a new inventory entry
	now := util.Timestamp()
	rsrc := podops.Resource{
		Name:       meta.Name,
		GUID:       meta.GUID,
		Kind:       podops.ResourceAsset,
		ParentGUID: production,
		Location:   location,
		Created:    now,
		Updated:    now,
	}

	if meta.IsImage() {
		rsrc.ImageURI = meta.Origin
		rsrc.ImageRel = rel
	} else {
		rsrc.EnclosureURI = meta.Origin
		rsrc.EnclosureRel = rel
	}

	if err := UpdateMetadata(ctx, meta); err != nil {
		return err
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

	if err := platform.DataStore().Delete(ctx, resourceKey(r.GUID)); err != nil {
		return err // FIXME put r back if this fails?
	}

	// validate the production after deleting a resource
	if err = ValidateProduction(ctx, prod); err != nil {
		p, err := GetProduction(ctx, prod)
		if err != nil {
			return err
		}
		p.BuildDate = 0 // FIXME BuildDate is the only flag we currently have to mark a production as VALID
		UpdateProduction(ctx, p)
	}

	if r.Kind == podops.ResourceAsset {
		if err := DeleteMetadata(ctx, guid); err != nil {
			return err
		}
		return RemoveAsset(ctx, prod, r.Location)
	}
	return RemoveResource(ctx, r.Location)
}

// ListResources returns all resources of type kind belonging to parentID
func ListResources(ctx context.Context, production, kind string) ([]*podops.Resource, error) {
	var r []*podops.Resource

	_kind, err := NormalizeKind(kind)
	if err != nil {
		return nil, err
	}

	if _kind == podops.ResourceALL {
		if _, err := platform.DataStore().GetAll(ctx, datastore.NewQuery(DatastoreResources).Filter("ParentGUID =", production).Order("-Created"), &r); err != nil {
			return nil, err
		}
	} else if _kind == podops.ResourceShow {
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

	if r.Kind == podops.ResourceAsset {
		meta, err := GetMetadata(ctx, guid)
		if err != nil {
			return nil, err
		}

		asset := podops.Asset{
			URI:   r.GetPublicLocation(),
			Title: r.Name,
			Type:  meta.ContentType,
			Size:  int(meta.Size),
			Rel:   podops.ResourceTypeLocal,
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

	bkt := platform.Storage().Bucket(podops.BucketProduction)
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

	bkt := platform.Storage().Bucket(podops.BucketProduction)
	reader, err := bkt.Object(path).NewReader(ctx)
	if err != nil {
		return nil, "", "", err
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, "", "", err
	}

	return loader.LoadResource(data)
}

// RemoveResource removes a resource from Cloud Storage
func RemoveResource(ctx context.Context, location string) error {
	bkt := platform.Storage().Bucket(podops.BucketProduction)

	obj := bkt.Object(location)
	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return errordef.ErrNoSuchResource
	}

	return bkt.Object(location).Delete(ctx)
}

// RemoveAsset removes a asset from Cloud Storage
func RemoveAsset(ctx context.Context, prod, location string) error {

	uri := fmt.Sprintf("%s/%s?l=%s", syncTaskEndpoint, prod, url.QueryEscape(location))
	// dispatch a request for background deletion
	_, err := p.CreateHttpTask(ctx, tasks.HttpMethod_DELETE, uri, env.GetString("PODOPS_API_KEY", ""), nil)

	return err
}

// UpdateShow is a helper function to update a show resource
func UpdateShow(ctx context.Context, location string, show *podops.Show) error {
	r, _ := GetResource(ctx, show.GUID())

	if r != nil {
		// resource already exists, just update the inventory
		if r.Kind != show.Kind {
			return fmt.Errorf("can not modify resource: expected '%s', received '%s'", r.Kind, show.Kind)
		}
		r.Name = show.Metadata.Name
		r.Location = location
		r.Title = show.Description.Title
		r.Summary = show.Description.Summary
		r.ImageURI = show.Image.ResolveURI(podops.DefaultStorageEndpoint, show.GUID())
		r.ImageRel = show.Image.Rel
		r.Updated = util.Timestamp()

		return updateResource(ctx, r)
	}

	// create a new inventory entry
	now := util.Timestamp()
	rsrc := podops.Resource{
		Name:       show.Metadata.Name,
		GUID:       show.GUID(),
		Kind:       podops.ResourceShow,
		ParentGUID: show.GUID(),
		Location:   location,
		Title:      show.Description.Title,
		Summary:    show.Description.Summary,
		ImageURI:   show.Image.ResolveURI(podops.DefaultStorageEndpoint, show.GUID()),
		ImageRel:   show.Image.Rel,
		Created:    now,
		Updated:    now,
	}
	return updateResource(ctx, &rsrc)
}

// UpdateEpisode is a helper function to update a episode resource
func UpdateEpisode(ctx context.Context, location string, episode *podops.Episode) error {
	// check if resource with same name already exists for the parent production
	rn, err := FindResource(ctx, episode.Parent(), episode.Metadata.Name)
	if err != nil {
		return err
	}
	r, err := GetResource(ctx, episode.GUID())
	if err != nil {
		return err
	}

	if rn != nil && r != nil {
		if rn.GUID != r.GUID {
			return fmt.Errorf("can not update resource: '%s/%s' already exists", episode.Parent(), episode.Metadata.Name)
		}
	}

	if r != nil {
		// resource already exists, just update the inventory
		if r.Kind != episode.Kind {
			return fmt.Errorf("can not modify resource: expected '%s', received '%s'", r.Kind, episode.Kind)
		}
		index, _ := strconv.ParseInt(episode.Metadata.Labels[podops.LabelEpisode], 10, 64)

		r.Name = episode.Metadata.Name
		r.ParentGUID = episode.Metadata.Labels[podops.LabelParentGUID]
		r.Location = location
		r.Title = episode.Description.Title
		r.Summary = episode.Description.Summary
		r.Published = episode.PublishDateTimestamp()
		r.Index = int(index) // episode number
		r.EnclosureURI = episode.Enclosure.ResolveURI(podops.DefaultStorageEndpoint, episode.Parent())
		r.EnclosureRel = episode.Enclosure.Rel
		r.ImageURI = episode.Image.ResolveURI(podops.DefaultStorageEndpoint, episode.Parent())
		r.ImageRel = episode.Image.Rel
		r.Updated = util.Timestamp()

		return updateResource(ctx, r)
	}

	// create a new inventory entry
	now := util.Timestamp()
	index, _ := strconv.ParseInt(episode.Metadata.Labels[podops.LabelEpisode], 10, 64)

	rsrc := podops.Resource{
		Name:         episode.Metadata.Name,
		GUID:         episode.GUID(),
		Kind:         podops.ResourceEpisode,
		ParentGUID:   episode.Metadata.Labels[podops.LabelParentGUID],
		Location:     location,
		Title:        episode.Description.Title,
		Summary:      episode.Description.Summary,
		Published:    episode.PublishDateTimestamp(),
		Index:        int(index), // episode number
		EnclosureURI: episode.Enclosure.ResolveURI(podops.DefaultStorageEndpoint, episode.Parent()),
		EnclosureRel: episode.Enclosure.Rel,
		ImageURI:     episode.Image.ResolveURI(podops.DefaultStorageEndpoint, episode.Parent()),
		ImageRel:     episode.Image.Rel,
		Created:      now,
		Updated:      now,
	}
	return updateResource(ctx, &rsrc)
}

// EnsureAsset validates the existence of the asset and imports it if necessary
func EnsureAsset(ctx context.Context, production string, rsrc *podops.Asset) error {
	if rsrc.Rel == podops.ResourceTypeExternal {
		_, err := pingURL(rsrc.URI)
		return err
	}
	if rsrc.Rel == podops.ResourceTypeLocal {
		// FIXME replace later with checking of the ResourceMetadata entries ...
		path := fmt.Sprintf("%s/%s/%s", podops.DefaultStorageEndpoint, production, rsrc.URI)
		_, err := pingURL(path) // ping the CDN
		if err != nil {
			return err
		}
		return nil
	}
	if rsrc.Rel == podops.ResourceTypeImport {
		_, err := pingURL(rsrc.URI) // ping the URL already here to avoid queueing a request that will fail later anyways
		if err != nil {
			return err
		}

		// FIXME compare to ResourceMetadata first ...

		// dispatch a request for background import
		ir := podops.ImportRequest{
			GUID:     production,
			Source:   rsrc.URI,
			Original: rsrc.AssetName(),
		}
		_, err = p.CreateHttpTask(ctx, tasks.HttpMethod_POST, importTaskEndpoint, env.GetString("PODOPS_API_KEY", ""), &ir)
		if err != nil {
			return err
		}
	}
	return nil
}

// pingURL tries a HEAD or GET request to verify that 'url' exists and is reachable
func pingURL(url string) (http.Header, error) {

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", podops.UserAgentString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		defer resp.Body.Close()
		// anything other than OK, Created, Accepted, NoContent is treated as an error
		if resp.StatusCode > http.StatusNoContent {
			return nil, fmt.Errorf("can not verify '%s'", url)
		}
	}
	return resp.Header.Clone(), nil
}

// updateResource does what the name suggests
func updateResource(ctx context.Context, r *podops.Resource) error {
	if _, err := platform.DataStore().Put(ctx, resourceKey(r.GUID), r); err != nil {
		return err
	}
	return nil
}

func resourceKey(guid string) *datastore.Key {
	return datastore.NameKey(DatastoreResources, guid, nil)
}
