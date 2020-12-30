package resources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/platform/pkg/platform"

	"github.com/podops/podops/internal/config"
	t "github.com/podops/podops/internal/types"
)

const (
	// ImportTask route to ImportTaskEndpoint
	ImportTask = "/import"

	// full canonical route
	importTaskWithPrefix = "/_t/import"
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

// ImportTaskEndpoint implements async file import
func ImportTaskEndpoint(c *gin.Context) {
	var req t.ImportRequest

	err := c.BindJSON(&req)
	if err != nil {
		// just report and return, resending will not change anything
		platform.ReportError(err)
		c.Status(http.StatusOK)
		return
	}

	status := importResource(appengine.NewContext(c.Request), req.Source, req.Dest)
	c.Status(status)
}

func importResource(ctx context.Context, src, dest string) int {
	resp, err := http.Get(src)
	if err != nil {
		return resp.StatusCode
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		platform.ReportError(fmt.Errorf("Can not retrieve '%s': %s", src, resp.Status))
		return http.StatusBadRequest
	}

	meta := ExtractMetadataFromResponse(resp)
	obj := platform.Storage().Bucket(config.BucketCDN).Object(dest)
	writer := obj.NewWriter(ctx)
	writer.ContentType = meta.ContentType
	defer writer.Close()

	// transfer using a buffer
	buffer := make([]byte, 32768)
	l, err := io.CopyBuffer(writer, resp.Body, buffer)

	// error handling & verification
	if err != nil {
		platform.ReportError(fmt.Errorf("Can not transfer '%s': %v", dest, err))
		return http.StatusBadRequest
	}
	if l != meta.Size {
		platform.ReportError(fmt.Errorf("Error transfering '%s': Expected %d, reveived %d", src, meta.Size, l))
		return http.StatusBadRequest
	}

	// FIXME write metadata

	return http.StatusOK
}

/*
func importResource(ctx context.Context, src, dest string) int {
	resp, err := http.Get(src)
	if err != nil {
		return resp.StatusCode
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		platform.ReportError(fmt.Errorf("Can not retrieve '%s': %s", src, resp.Status))
		return http.StatusBadRequest
	}

	// FIXME this might not work for large files
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		platform.ReportError(fmt.Errorf("Can not retrieve '%s': %v", src, err))
		return http.StatusBadRequest
	}

	obj := platform.Storage().Bucket(config.BucketCDN).Object(dest)
	writer := obj.NewWriter(ctx)
	defer writer.Close()

	if _, err := writer.Write(data); err != nil {
		platform.ReportError(fmt.Errorf("Can not write '%s': %v", dest, err))
		return http.StatusBadRequest
	}

	return http.StatusOK
}
*/

// ExtractMetadataFromResponse extracts the metadata from http.Response
func ExtractMetadataFromResponse(resp *http.Response) *ContentMetadata {
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
