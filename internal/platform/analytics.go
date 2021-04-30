package platform

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	p "github.com/txsvc/platform"
	"github.com/txsvc/platform/pkg/env"
	"github.com/txsvc/platform/pkg/id"
)

const (
	analyticsEndpoint = "https://www.google-analytics.com/collect"

	typePageView    = "pageview"
	typeScreenView  = "screenview"
	typeEvent       = "event"
	typeTransaction = "transaction"
	typeItem        = "item"
	typeSocial      = "social"
)

type (
	// Event  contains, events
	Event struct {
		Category string
		Action   string
		Label    string
		Value    int
	}
)

var (
	// measurementID is the Google Analytics ID
	measurementID string = env.GetString("MEASUREMENT_ID", "UA-xxxxxxxx-x")
	appID         string = env.GetString("SERVICE_NAME", "backend")
	filterPages   []string
	events        chan (*url.Values)
)

func init() {
	filterPages = make([]string, 3)
	filterPages[0] = "/assets/js"
	filterPages[1] = "/assets/css"
	filterPages[2] = "/assets/static"

	events = make(chan *url.Values, 100)
	go upload()
}

// PageViewMiddleware logs page views to Google Analytics
func PageViewMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}

		// skip /assets/..
		path := c.Request().URL.Path
		for _, s := range filterPages {
			if strings.HasPrefix(path, s) {
				return nil
			}
		}

		v := make(map[string]string)
		v["t"] = typePageView
		v["ds"] = "web"

		if err := PostToAnalytics(c.Request(), &v); err != nil {
			p.ReportError(err)
			return err
		}

		return nil
	}
}

// TrackEvent post an event to analytics
func TrackEvent(request *http.Request, category, action, label string, value int) error {
	v := make(map[string]string)

	v["t"] = typeEvent
	v["ec"] = category
	v["ea"] = action
	v["el"] = label
	v["ev"] = strconv.FormatInt(int64(value), 10)
	v["ds"] = appID

	if err := PostToAnalytics(request, &v); err != nil {
		p.ReportError(err)
		return err
	}
	return nil
}

// PostToAnalytics send the values to Google Analytics
func PostToAnalytics(request *http.Request, values *map[string]string) error {

	ip := request.RemoteAddr
	userAgent := request.UserAgent()
	uid := id.Fingerprint(userAgent + ip)
	path := request.URL.Path
	dl := request.Host + request.RequestURI

	// the basics
	formValues := url.Values{
		"v":   {"1"},
		"tid": {measurementID},
		"uid": {uid},
		"uip": {ip},

		"dl":  {url.QueryEscape(dl)},
		"ua":  {url.QueryEscape(userAgent)},
		"dh":  {url.QueryEscape(request.Host)},
		"dp":  {url.QueryEscape(path)},
		"dt":  {path},
		"npa": {"1"}, // Disabling Advertising Personalization

	}

	// event specific k/v
	for k, v := range *values {
		vv := make([]string, 1)
		vv[0] = url.QueryEscape(v)
		formValues[k] = vv
	}

	events <- &formValues
	/*
		resp, err := http.PostForm(analyticsEndpoint, formValues)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			return fmt.Errorf("Google Analytics returned '%d'", resp.StatusCode)
		}
	*/
	return nil
}

func upload() {
	for {
		v := <-events

		resp, err := http.PostForm(analyticsEndpoint, *v)
		if err == nil && resp != nil {
			resp.Body.Close()
		} else {
			p.ReportError(err)
		}
	}
}
