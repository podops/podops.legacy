package cdn

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fupas/commons/pkg/util"
	"github.com/labstack/echo/v4"

	"github.com/podops/podops"
	"github.com/podops/podops/apiv1"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/internal/platform"
)

// ImportTaskEndpoint implements async file import from a remote source into the CDN
func ImportTaskEndpoint(c echo.Context) error {
	var req podops.ImportRequest

	err := c.Bind(&req)
	if err != nil {
		// just report and return, resending will not change anything
		platform.ReportError(err)
		return c.NoContent(http.StatusOK)
	}

	if req.GUID == "" || req.Source == "" || req.Original == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	ctx := platform.NewHttpContext(c)

	if err := apiv1.AuthorizeAccessProduction(ctx, c, apiv1.ScopeAPIAdmin, req.GUID); err != nil {
		return platform.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	status := ImportResource(ctx, req.GUID, req.Source, req.Original)
	return c.NoContent(status)
}

// ImportResource imports a resource from src and places it into the CDN
func ImportResource(ctx context.Context, prod, src, original string) int {
	resp, err := http.Get(src)
	if err != nil {
		return resp.StatusCode
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		platform.ReportError(fmt.Errorf("can not retrieve '%s': %s", src, resp.Status))
		return http.StatusBadRequest
	}

	// update the inventory
	meta := extractMetadataFromResponse(resp)

	temp := podops.Asset{
		URI: src,
		Rel: podops.ResourceTypeImport,
	}
	parts := strings.Split(temp.FingerprintURI(prod), "/")

	meta.Name = parts[len(parts)-1:][0]
	meta.GUID = util.Checksum(src)

	relPath := prod + "/" + meta.Name
	path := filepath.Join(podops.StorageLocation, relPath)

	// FIXME check metadata and avoid downloading if still valid?

	os.MkdirAll(filepath.Dir(path), os.ModePerm) // make sure sub-folders exist
	out, err := os.Create(path)
	if err != nil {
		platform.ReportError(fmt.Errorf("can not transfer '%s': %v", src, err))
		return http.StatusBadRequest
	}
	defer out.Close()

	// transfer using a buffer
	buffer := make([]byte, 65536)
	l, err := io.CopyBuffer(out, resp.Body, buffer)

	// error handling & verification
	if err != nil {
		platform.ReportError(fmt.Errorf("can not transfer '%s': %v", src, err))
		return http.StatusBadRequest
	}
	if l != meta.Size {
		platform.ReportError(fmt.Errorf("error transfering '%s': expected %d, reveived %d", src, meta.Size, l))
		return http.StatusBadRequest
	}

	// explicitly close the file here
	out.Close()

	//
	meta.Duration = calculateLength(meta.ContentType, path)

	// FIXME write metadata ?

	if err := backend.UpdateAsset(ctx, meta.Name, meta.GUID, podops.ResourceAsset, prod, relPath, meta.ContentType, original, meta.Etag, meta.Size, meta.Duration); err != nil {
		platform.ReportError(fmt.Errorf("error updating inventory: %v", err))
		return http.StatusBadRequest
	}

	return http.StatusOK
}

// extractMetadataFromResponse extracts the metadata from http.Response
func extractMetadataFromResponse(resp *http.Response) *podops.ContentMetadata {
	meta := podops.ContentMetadata{
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

// calculateLength returns the play duration of a media file like a .mp3
func calculateLength(contentType, path string) int64 {
	return 0 // FIXME to be implemented
}
