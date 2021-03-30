package backend

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"

	"github.com/fupas/commons/pkg/util"
	ds "github.com/fupas/platform/pkg/platform"

	a "github.com/podops/podops"
	"github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/platform"
)

type (
	// ContentMetadata keeps basic data on resource
	ContentMetadata struct {
		Size        int64
		Duration    int64
		ContentType string
		Etag        string
		Timestamp   int64
	}
)

// ImportResource import a resource from a src and place it into the CDN
func ImportResource(ctx context.Context, src, dest, original string) int {
	resp, err := http.Get(src)
	if err != nil {
		return resp.StatusCode
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		platform.ReportError(fmt.Errorf("can not retrieve '%s': %s", src, resp.Status))
		return http.StatusBadRequest
	}

	meta := extractMetadataFromResponse(resp)
	obj := ds.Storage().Bucket(a.BucketCDN).Object(dest)
	writer := obj.NewWriter(ctx)
	writer.ContentType = meta.ContentType
	defer writer.Close()

	// transfer using a buffer
	buffer := make([]byte, 65536)
	l, err := io.CopyBuffer(writer, resp.Body, buffer)

	// error handling & verification
	if err != nil {
		platform.ReportError(fmt.Errorf("can not transfer '%s': %v", dest, err))
		return http.StatusBadRequest
	}
	if l != meta.Size {
		platform.ReportError(fmt.Errorf("error transfering '%s': expected %d, reveived %d", src, meta.Size, l))
		return http.StatusBadRequest
	}

	// update the inventory
	parent := strings.Split(dest, "/")[0]

	temp := a.Asset{
		URI: src,
		Rel: a.ResourceTypeImport,
	}

	name := strings.Split(temp.FingerprintURI(parent), "/")[1]
	duration := int64(0) // FIXME implement it

	if err := UpdateAsset(ctx, name, util.Checksum(src), a.ResourceAsset, parent, dest, meta.ContentType, original, meta.Etag, meta.Size, duration); err != nil {
		platform.ReportError(fmt.Errorf("error updating inventory: %v", err))
		return http.StatusBadRequest
	}

	return http.StatusOK
}

// EnsureAsset validates the existence of the asset and imports it if necessary
func EnsureAsset(ctx context.Context, production string, rsrc *a.Asset) error {
	if rsrc.Rel == a.ResourceTypeExternal {
		_, err := pingURL(rsrc.URI)
		return err
	}
	if rsrc.Rel == a.ResourceTypeLocal {
		path := fmt.Sprintf("%s/%s", production, rsrc.URI)
		if !resourceExists(ctx, path) {
			return fmt.Errorf("can not find '%s'", rsrc.URI)
		}
		return nil
	}
	if rsrc.Rel == a.ResourceTypeImport {
		_, err := pingURL(rsrc.URI) // ping the URL already here to avoid queueing a request that will fail later anyways
		if err != nil {
			return err
		}

		path := rsrc.FingerprintURI(production)
		if resourceExists(ctx, path) { // do nothing as the asset is present FIXME re-download if --force is set
			return nil // FIXME verify that the asset is unchanged, otherwise re-import
		}

		// dispatch a request for background import
		_, err = platform.CreateTask(ctx, apiv1.ImportTaskWithPrefix, &a.ImportRequest{Source: rsrc.URI, Dest: path, Original: rsrc.AssetName()})
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
	req.Header.Set("User-Agent", a.UserAgentString)

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
	obj := ds.Storage().Bucket(a.BucketCDN).Object(path)
	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false
	}
	return true
}

// extractMetadataFromResponse extracts the metadata from http.Response
func extractMetadataFromResponse(resp *http.Response) *ContentMetadata {
	meta := ContentMetadata{
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
