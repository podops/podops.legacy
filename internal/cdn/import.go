package cdn

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/pkg/api"
	"github.com/txsvc/platform/v2/pkg/authentication"

	"github.com/podops/podops"
	"github.com/podops/podops/apiv1"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/internal/messagedef"
	"github.com/podops/podops/internal/metadata"
)

// ImportTaskEndpoint implements async file import from a remote source into the CDN
func ImportTaskEndpoint(c echo.Context) error {
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

	if err := apiv1.AuthorizeAccessProduction(ctx, c, authentication.ScopeAPIAdmin, req.GUID); err != nil {
		return api.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	status := ImportResource(ctx, req.GUID, req.Source)
	return c.NoContent(status)
}

// ImportResource imports a resource from src and places it into the CDN
func ImportResource(ctx context.Context, prod, src string) int {
	resp, err := http.Get(src)
	if err != nil {
		return resp.StatusCode
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		platform.ReportError(fmt.Errorf(messagedef.MsgResourceImportError, src))
		return http.StatusBadRequest
	}

	// update the inventory
	meta := metadata.ExtractMetadataFromResponse(resp)

	meta.Name = metadata.LocalNamePart(metadata.FingerprintWithExt(prod, src))
	meta.Origin = src
	meta.GUID = metadata.FingerprintURI(prod, src)
	meta.ParentGUID = prod

	relPath := prod + "/" + meta.Name
	path := filepath.Join(podops.StorageLocation, relPath)

	// FIXME check metadata and avoid downloading if still valid?

	os.MkdirAll(filepath.Dir(path), os.ModePerm) // make sure sub-folders exist
	out, err := os.Create(path)
	if err != nil {
		platform.ReportError(err)
		return http.StatusBadRequest
	}
	defer out.Close()

	// transfer using a buffer
	buffer := make([]byte, 65536)
	l, err := io.CopyBuffer(out, resp.Body, buffer)

	// error handling & verification
	if err != nil {
		platform.ReportError(err)
		return http.StatusBadRequest
	}
	if l != meta.Size {
		platform.ReportError(fmt.Errorf(messagedef.MsgResourceImportError, src))
		return http.StatusBadRequest
	}

	// explicitly close the file here
	out.Close()

	// calculate the length of an audio file, if it is an audio file
	if meta.IsAudio() {
		meta.Duration, _ = metadata.CalculateLength(path)
	}

	// update the inventory
	if err := backend.UpdateAsset(ctx, meta, prod, relPath, podops.ResourceTypeImport); err != nil {
		platform.ReportError(err)
		return http.StatusBadRequest
	}

	return http.StatusOK
}
