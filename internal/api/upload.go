package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/platform"
	"github.com/labstack/echo/v4"

	a "github.com/podops/podops/apiv1"
	p "github.com/podops/podops/internal/platform"
	"github.com/podops/podops/pkg/api"
	"github.com/podops/podops/pkg/backend"
)

// UploadEndpoint implements content upload
func UploadEndpoint(c echo.Context) error {
	ctx := api.NewHttpContext(c)

	mr, err := c.Request().MultipartReader()
	if err != nil {
		return api.ErrorResponse(c, http.StatusInternalServerError, err)
	}
	prod := c.Param("prod")
	if prod == "" {
		return api.ErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid route, expected ':prod'"))
	}

	if err := AuthorizeAccessProduction(ctx, c, scopeResourceWrite, prod); err != nil {
		return api.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return api.ErrorResponse(c, http.StatusInternalServerError, err)
		}

		if p.FormName() == "asset" {
			location := fmt.Sprintf("%s/%s", prod, p.FileName())

			bkt := platform.Storage().Bucket(a.BucketCDN)
			obj := bkt.Object(location)
			writer := obj.NewWriter(ctx)
			defer writer.Close() // just to be sure we really close the writer

			if _, err := io.Copy(writer, p); err != nil {
				return api.ErrorResponse(c, http.StatusInternalServerError, err)
			}
			writer.Close() // force close to have attributes like size etc correct

			// get the attributes back
			attr, err := obj.Attrs(ctx)
			if err != nil {
				return api.ErrorResponse(c, http.StatusInternalServerError, err)
			}

			duration := int64(0) // FIXME implement it
			original := p.FileName()

			// update the inventory
			backend.UpdateAsset(ctx, p.FileName(), util.Checksum(location), a.ResourceAsset, prod, location, attr.ContentType, original, attr.Etag, attr.Size, duration)
		}
	}

	// track api access for billing etc
	p.TrackEvent(c.Request(), "api", "upload", prod, 1)

	return c.NoContent(http.StatusCreated)
}
