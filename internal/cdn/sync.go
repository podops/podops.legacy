package cdn

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	cs "github.com/fupas/platform/pkg/platform"
	"github.com/labstack/echo/v4"

	"github.com/podops/podops"
	"github.com/podops/podops/apiv1"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/internal/errordef"
	"github.com/podops/podops/internal/platform"
)

// SyncTaskEndpoint syncs files between the cloud storage and the CDN
func SyncTaskEndpoint(c echo.Context) error {
	var req podops.ImportRequest

	err := c.Bind(&req)
	if err != nil {
		// just report and return, resending will not change anything
		platform.ReportError(err)
		return c.NoContent(http.StatusOK)
	}

	if req.GUID == "" || req.Source == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	ctx := platform.NewHttpContext(c)

	if err := apiv1.AuthorizeAccessProduction(ctx, c, apiv1.ScopeAPIAdmin, req.GUID); err != nil {
		return platform.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	status := SyncResource(ctx, req.GUID, req.Source)
	return c.NoContent(status)
}

// DeleteTaskEndpoint removes files from the CDN
func DeleteTaskEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c)

	prod := c.Param("prod")
	kind := c.Param("kind")
	guid := c.Param("id")

	if !apiv1.ValidateNotEmpty(prod, kind, guid) {
		return platform.ErrorResponse(c, http.StatusBadRequest, errordef.ErrInvalidRoute)
	}
	if err := apiv1.AuthorizeAccessResource(ctx, c, apiv1.ScopeAPIAdmin, guid); err != nil {
		return platform.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	r, err := backend.GetResource(ctx, guid)
	if err != nil {
		return platform.ErrorResponse(c, http.StatusBadRequest, err)
	}

	status := DeleteResource(ctx, prod, r.Location)
	return c.NoContent(status)
}

// SyncResource imports a resource from the cloud storage and places it into the CDN
func SyncResource(ctx context.Context, prod, src string) int {
	relPath := prod + "/" + src

	bkt := cs.Storage().Bucket(podops.BucketProduction)
	reader, err := bkt.Object(relPath).NewReader(ctx)
	if err != nil {
		platform.ReportError(fmt.Errorf("can not transfer '%s': %v", src, err))
		return http.StatusBadRequest
	}

	path := filepath.Join(podops.StorageLocation, relPath)

	os.MkdirAll(filepath.Dir(path), os.ModePerm) // make sure sub-folders exist
	out, err := os.Create(path)
	if err != nil {
		platform.ReportError(fmt.Errorf("can not transfer '%s': %v", src, err))
		return http.StatusBadRequest
	}
	defer out.Close()

	// transfer the file
	_, err = io.Copy(out, reader)
	if err != nil {
		platform.ReportError(fmt.Errorf("can not transfer '%s': %v", src, err))
		return http.StatusBadRequest
	}

	return http.StatusOK
}

// DeleteResource removes a resource from the CDN
func DeleteResource(ctx context.Context, prod, location string) int {
	path := filepath.Join(podops.StorageLocation, location)
	err := os.Remove(path)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
