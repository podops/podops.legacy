package cdn

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fupas/commons/pkg/util"
	"github.com/labstack/echo/v4"

	"github.com/podops/podops"
	"github.com/podops/podops/apiv1"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/internal/metadata"
	p "github.com/podops/podops/internal/platform"
)

// UploadEndpoint implements content upload
func UploadEndpoint(c echo.Context) error {
	ctx := p.NewHttpContext(c)

	mr, err := c.Request().MultipartReader()
	if err != nil {
		return p.ErrorResponse(c, http.StatusInternalServerError, err)
	}
	prod := c.Param("prod")
	if prod == "" {
		return p.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':prod'"))
	}

	if err := apiv1.AuthorizeAccessProduction(ctx, c, apiv1.ScopeResourceWrite, prod); err != nil {
		return p.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return p.ErrorResponse(c, http.StatusInternalServerError, err)
		}

		if part.FormName() == "asset" {
			location := fmt.Sprintf("%s/%s", prod, part.FileName())
			path := filepath.Join(podops.StorageLocation, location)

			os.MkdirAll(filepath.Dir(path), os.ModePerm) // make sure sub-folders exist
			out, err := os.Create(path)
			if err != nil {
				return p.ErrorResponse(c, http.StatusInternalServerError, err)
			}
			defer out.Close()

			if _, err := io.Copy(out, part); err != nil {
				return p.ErrorResponse(c, http.StatusInternalServerError, err)
			}
			out.Close() // force close to have attributes like size etc correct

			// extract the metadata from the file
			meta, err := metadata.ExtractMetadataFromFile(path)
			if err != nil {
				return p.ErrorResponse(c, http.StatusInternalServerError, err)
			}
			meta.GUID = util.Checksum(meta.Name)

			// update the inventory
			if err := backend.UpdateResourceMetadata(ctx, meta); err != nil {
				return p.ErrorResponse(c, http.StatusInternalServerError, err)
			}
			if err := backend.UpdateAsset(ctx, meta.Name, meta.GUID, podops.ResourceAsset, prod, location, meta.ContentType, meta.Name, meta.Etag, meta.Size, meta.Duration); err != nil {
				return p.ErrorResponse(c, http.StatusInternalServerError, err)
			}
		}
	}

	// track api access for billing etc
	p.TrackEvent(c.Request(), "api", "upload", prod, 1)

	return c.NoContent(http.StatusCreated)
}
