package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/podops/podops/internal/platform"
	"github.com/podops/podops/pkg/api"
	"github.com/podops/podops/pkg/backend"
	"github.com/podops/podops/pkg/backend/models"
)

// ImportTaskEndpoint implements async file import
func ImportTaskEndpoint(c echo.Context) error {
	var req *models.ImportRequest = new(models.ImportRequest)

	err := c.Bind(req)
	if err != nil {
		// just report and return, resending will not change anything
		platform.ReportError(err)
		return c.NoContent(http.StatusOK)
	}

	// FIXME does it make sense to retry? If not, send StatusOK
	status := backend.ImportResource(api.NewHttpContext(c), req.Source, req.Dest, req.Original)
	return c.NoContent(status)
}
