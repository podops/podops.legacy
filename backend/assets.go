package backend

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/pkg/env"
	"github.com/txsvc/platform/v2/pkg/timestamp"
	"github.com/txsvc/platform/v2/tasks"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/messagedef"
	"github.com/podops/podops/internal/metadata"
	"github.com/podops/podops/internal/transport"
)

// UpdateAsset updates the resource inventory
func UpdateAsset(ctx context.Context, meta *metadata.Metadata, production, location, rel string) error {
	r, _ := GetResource(ctx, meta.GUID)

	if r != nil {
		// resource already exists, just update the inventory
		r.Name = meta.Name
		r.ParentGUID = production
		r.Location = location
		r.Updated = timestamp.Now()

		if meta.IsImage() {
			r.ImageURI = fmt.Sprintf("%s/%s", podops.DefaultStorageEndpoint, location)
			r.ImageRel = rel
		} else {
			r.EnclosureURI = fmt.Sprintf("%s/%s", podops.DefaultStorageEndpoint, location)
			r.EnclosureRel = rel
		}

		if err := UpdateMetadata(ctx, meta); err != nil {
			return err
		}
		return updateResource(ctx, r)
	}

	// create a new inventory entry
	now := timestamp.Now()
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
		rsrc.ImageURI = fmt.Sprintf("%s/%s", podops.DefaultStorageEndpoint, location)
		rsrc.ImageRel = rel
	} else {
		rsrc.EnclosureURI = fmt.Sprintf("%s/%s", podops.DefaultStorageEndpoint, location)
		rsrc.EnclosureRel = rel
	}

	if err := UpdateMetadata(ctx, meta); err != nil {
		return err
	}
	return updateResource(ctx, &rsrc)
}

// RemoveAsset removes a asset from Cloud Storage
func RemoveAsset(ctx context.Context, prod, location string) error {

	//uri := fmt.Sprintf("%s/%s?l=%s", syncTaskEndpoint, prod, url.QueryEscape(location))
	// dispatch a request for background deletion
	task := tasks.HttpTask{
		Method:  tasks.HttpMethodDelete,
		Request: fmt.Sprintf("%s/%s?l=%s", syncTaskEndpoint, prod, url.QueryEscape(location)),
		Token:   env.GetString("PODOPS_API_KEY", ""),
		Payload: nil,
	}
	err := platform.NewTask(task)

	//_, err := p.CreateHttpTask(ctx, tasks.HttpMethod_DELETE, uri, env.GetString("PODOPS_API_KEY", ""), nil)

	return err
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
		ir := podops.SyncRequest{
			GUID:   production,
			Source: rsrc.URI,
		}

		task := tasks.HttpTask{
			Method:  tasks.HttpMethodPost,
			Request: importTaskEndpoint,
			Token:   env.GetString("PODOPS_API_KEY", ""),
			Payload: &ir,
		}
		err = platform.NewTask(task)
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
	req.Header.Set("User-Agent", transport.UserAgentString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		defer resp.Body.Close()
		// anything other than OK, Created, Accepted, NoContent is treated as an error
		if resp.StatusCode > http.StatusNoContent {
			return nil, fmt.Errorf(messagedef.MsgResourceIsInvalid, url)
		}
	}
	return resp.Header.Clone(), nil
}
