package cdn

/* See the following resourced for reference:

https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests
https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers

https://github.com/gin-gonic/gin

https://cloud.google.com/cdn/docs/

*/

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/env"
	"github.com/txsvc/platform/pkg/platform"

	t "github.com/podops/podops/internal/types"
)

const (
	cacheControl = "public, max-age=1800"
)

type (
	// Header extracts the relevant HTTP header stuff
	Header struct {
		Range           string `header:"Range"`
		UserAgent       string `header:"User-Agent"`
		Forwarded       string `header:"Forwarded"`
		XForwardedFor   string `header:"X-Forwarded-For"`
		XForwwardedHost string `header:"X-Forwarded-Host"`
		Referer         string `header:"Referer"`
	}
)

var redirectBase string
var bkt *storage.BucketHandle

func init() {
	redirectBase = env.GetString("REDIRECT_URL", "https://storage.googleapis.com/cdn.podops.dev")

	bkt = platform.Storage().Bucket(t.BucketCDN)
}

// ServeContentEndpoint handles request for content
func ServeContentEndpoint(c *gin.Context) {

	// return an error if the request is anything other than GET/HEAD
	m := c.Request.Method
	if m != "" && m != "GET" && m != "HEAD" {
		platform.ReportError(fmt.Errorf("cdn: received a '%s' request", m))
		c.Status(http.StatusBadRequest)
		return
	}

	// extract headers we are interested in
	header := Header{}
	err := c.ShouldBindHeader(&header)
	if err != nil {
		platform.ReportError(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	ctx := appengine.NewContext(c.Request)

	// get attributes, can be cached ...
	rsrc := c.Request.URL.Path[1:]
	obj := bkt.Object(rsrc)
	attr, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		platform.ReportError(fmt.Errorf("cdn: can not find resource '%s'", rsrc))
		c.Status(http.StatusNotFound)
		return
	}

	// handle HEAD request
	if m == "HEAD" {
		c.Header("etag", attr.Etag)
		c.Header("accept-ranges", "bytes")
		c.Header("cache-control", cacheControl)
		c.Header("accept-ranges", "bytes")
		c.Header("content-type", attr.ContentType)
		c.Header("content-length", fmt.Sprintf("%d", attr.Size))

		c.Status(http.StatusOK)
		return
	}

	// GET from here on

	// FIXME this is just a hack!
	platform.Log(header, nil)

	// create a reader
	//reader, err := obj.NewReader(ctx)
	offset, length := parseRange(header.Range)
	reader, err := obj.NewRangeReader(ctx, offset, length)
	if err != nil {
		platform.ReportError(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// extra headers
	extraHeaders := map[string]string{
		"etag":          attr.Etag,
		"accept-ranges": "bytes",
		"cache-control": cacheControl,
	}

	// send the data back, no redirects

	status := http.StatusOK
	if length != -1 {
		status = http.StatusPartialContent
	}

	c.DataFromReader(status, attr.Size, attr.ContentType, reader, extraHeaders)

}

func parseRange(r string) (int64, int64) {
	if r == "" {
		return 0, -1 // no range requested
	}
	parts := strings.Split(r, "=")
	if len(parts) != 2 {
		return 0, -1 // no range requested
	}
	// we simply assume that parts[0] == "bytes"
	ra := strings.Split(parts[1], "-")
	if len(ra) != 2 { // again a simplification, multiple ranges or overlapping ranges are not supported
		return 0, -1
	}

	start, err := strconv.ParseInt(ra[0], 10, 64)
	if err != nil {
		return 0, -1
	}
	end, err := strconv.ParseInt(ra[1], 10, 64)
	if err != nil {
		return 0, -1
	}

	return start, end - start
}

/*

other headers:

Alt-Svc

*/
