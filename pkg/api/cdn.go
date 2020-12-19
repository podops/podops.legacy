package api

/* See the following resourced for reference:

https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests
https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers

https://github.com/gin-gonic/gin

https://cloud.google.com/cdn/docs/

*/

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/txsvc/commons/pkg/env"
	"github.com/txsvc/platform/pkg/platform"
)

var redirectBase string = env.GetString("REDIRECT_URL", "https://storage.googleapis.com/cdn.podops.dev")

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

// RedirectToStorageEndpoint redirects requests to Cloud Storage
func RedirectToStorageEndpoint(c *gin.Context) {

	// return an error if the request is anything other than GET/HEAD
	m := c.Request.Method
	if m != "" && m != "GET" && m != "HEAD" {
		c.Status(http.StatusBadRequest)
		return
	}

	// extract headers we are interested in
	header := Header{}
	err := c.ShouldBindHeader(&header)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// FIXME: implement analytics here ...
	target := redirectBase + c.Request.URL.Path
	platform.Log(header)

	// redirect to the CDN of Google's Cloud Storage
	c.Redirect(http.StatusTemporaryRedirect, target)
}
