package cdn

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/auth"
	"github.com/txsvc/platform/v2/pkg/api"
	ds "github.com/txsvc/platform/v2/pkg/datastore"
	"github.com/txsvc/platform/v2/pkg/validate"

	"github.com/podops/podops"
	"github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/errordef"
)

// SyncTaskEndpoint syncs files between the cloud storage and the CDN
func SyncTaskEndpoint(c echo.Context) error {
	var req podops.SyncRequest

	err := c.Bind(&req)
	if err != nil {
		// just report and return, resending will not change anything
		platform.ReportError(err)
		return c.NoContent(http.StatusOK)
	}

	if req.GUID == "" || req.Source == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	ctx := platform.NewHttpContext(c.Request())

	if err := apiv1.AuthorizeAccessProduction(ctx, c, auth.ScopeAPIAdmin, req.GUID); err != nil {
		return api.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	status := SyncResource(ctx, req.GUID, req.Source)
	return c.NoContent(status)
}

// DeleteTaskEndpoint removes a file from the CDN
func DeleteTaskEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c.Request())

	prod := c.Param("prod")
	location := c.QueryParam("l")

	if !validate.NotEmpty(prod, location) {
		return api.ErrorResponse(c, http.StatusBadRequest, errordef.ErrInvalidRoute)
	}
	if err := apiv1.AuthorizeAccessProduction(ctx, c, auth.ScopeAPIAdmin, prod); err != nil {
		// validate against production only, the resource is already gone by now
		return api.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	status := DeleteResource(ctx, location)
	return c.NoContent(status)
}

// SyncResource imports a resource from the cloud storage and places it into the CDN
func SyncResource(ctx context.Context, prod, src string) int {
	relPath := prod + "/" + src

	bkt := ds.Storage().Bucket(podops.BucketProduction)
	reader, err := bkt.Object(relPath).NewReader(ctx)
	if err != nil {
		platform.ReportError(err)
		return http.StatusBadRequest
	}

	path := filepath.Join(podops.StorageLocation, relPath)

	os.MkdirAll(filepath.Dir(path), os.ModePerm) // make sure sub-folders exist
	out, err := os.Create(path)
	if err != nil {
		platform.ReportError(err)
		return http.StatusBadRequest
	}
	defer out.Close()

	// transfer the file
	_, err = io.Copy(out, reader)
	if err != nil {
		platform.ReportError(err)
		return http.StatusBadRequest
	}

	return http.StatusOK
}

// DeleteResource removes a resource from the CDN
func DeleteResource(ctx context.Context, location string) int {
	path := filepath.Join(podops.StorageLocation, location)
	err := os.Remove(path)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
