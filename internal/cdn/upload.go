package cdn

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/txsvc/platform"
	"github.com/txsvc/platform/pkg/server"

	"github.com/podops/podops"
	"github.com/podops/podops/apiv1"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/internal/errordef"
	"github.com/podops/podops/internal/metadata"

	lp "github.com/podops/podops/internal/platform"
)

// UploadEndpoint implements content upload
func UploadEndpoint(c echo.Context) error {
	ctx := platform.NewHttpContext(c.Request())

	mr, err := c.Request().MultipartReader()
	if err != nil {
		return server.ErrorResponse(c, http.StatusInternalServerError, err)
	}
	prod := c.Param("prod")
	if prod == "" {
		return server.ErrorResponse(c, http.StatusBadRequest, errordef.ErrInvalidRoute)
	}

	if err := apiv1.AuthorizeAccessProduction(ctx, c, apiv1.ScopeResourceWrite, prod); err != nil {
		return server.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return server.ErrorResponse(c, http.StatusInternalServerError, err)
		}

		if part.FormName() == "asset" {
			location := fmt.Sprintf("%s/%s", prod, part.FileName())
			path := filepath.Join(podops.StorageLocation, location)

			os.MkdirAll(filepath.Dir(path), os.ModePerm) // make sure sub-folders exist
			out, err := os.Create(path)
			if err != nil {
				return server.ErrorResponse(c, http.StatusInternalServerError, err)
			}
			defer out.Close()

			if _, err := io.Copy(out, part); err != nil {
				return server.ErrorResponse(c, http.StatusInternalServerError, err)
			}
			out.Close() // force close to have attributes like size etc correct

			// extract the metadata from the file
			meta, err := metadata.ExtractMetadataFromFile(path)
			if err != nil {
				return server.ErrorResponse(c, http.StatusInternalServerError, err)
			}
			meta.GUID = metadata.FingerprintURI(prod, meta.Name)
			meta.ParentGUID = prod
			meta.Origin = location

			// update the inventory
			if err := backend.UpdateAsset(ctx, meta, prod, location, podops.ResourceTypeLocal); err != nil {
				return server.ErrorResponse(c, http.StatusInternalServerError, err)
			}
		}
	}

	// track api access for billing etc
	lp.TrackEvent(c.Request(), "api", "upload", prod, 1)

	return c.NoContent(http.StatusCreated)
}
