package backend

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/storage"

	ds "github.com/fupas/platform/pkg/platform"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/platform"
)

const (
	// full canonical route
	ImportTaskWithPrefix = "/_t/import"
)

// EnsureAsset validates the existence of the asset and imports it if necessary
func EnsureAsset(ctx context.Context, production string, rsrc *podops.Asset) error {
	if rsrc.Rel == podops.ResourceTypeExternal {
		_, err := pingURL(rsrc.URI)
		return err
	}
	if rsrc.Rel == podops.ResourceTypeLocal {
		path := fmt.Sprintf("%s/%s", production, rsrc.URI)
		if !resourceExists(ctx, path) {
			return fmt.Errorf("can not find '%s'", rsrc.URI)
		}
		return nil
	}
	if rsrc.Rel == podops.ResourceTypeImport {
		_, err := pingURL(rsrc.URI) // ping the URL already here to avoid queueing a request that will fail later anyways
		if err != nil {
			return err
		}

		path := rsrc.FingerprintURI(production)
		if resourceExists(ctx, path) { // do nothing as the asset is present FIXME re-download if --force is set
			return nil // FIXME verify that the asset is unchanged, otherwise re-import
		}

		// dispatch a request for background import
		_, err = platform.CreateTask(ctx, ImportTaskWithPrefix, &podops.ImportRequest{Source: rsrc.URI, Dest: path, Original: rsrc.AssetName()})
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

// resourceExists verifies the resource .yaml exists
func resourceExists(ctx context.Context, path string) bool {
	obj := ds.Storage().Bucket(podops.BucketCDN).Object(path)
	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false
	}
	return true
}
