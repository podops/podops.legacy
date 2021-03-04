package analytics

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/fupas/commons/pkg/env"
	"github.com/fupas/commons/pkg/util"
	"github.com/podops/podops/internal/observer"
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
)

// TrackEvent post an event to analytics
func TrackEvent(request *http.Request, category, action, label string, value int) error {
	v := make(map[string]string)

	v["t"] = typeEvent
	v["ec"] = category
	v["ea"] = action
	v["el"] = label
	v["ev"] = strconv.FormatInt(int64(value), 10)

	if err := PostToAnalytics(request, &v); err != nil {
		observer.ReportError(err)
		return err
	}
	return nil
}

// PostToAnalytics send the values to Google Analytics
func PostToAnalytics(request *http.Request, values *map[string]string) error {

	ip := request.RemoteAddr
	userAgent := request.UserAgent()
	uid := util.Fingerprint(userAgent + ip)

	// the basics
	formValues := url.Values{
		"v":   {"1"},
		"tid": {measurementID},
		"ds":  {appID},
		"uid": {uid},
		"uip": {ip},

		"ua":  {url.QueryEscape(userAgent)},
		"dh":  {url.QueryEscape(request.Host)},
		"dp":  {url.QueryEscape(request.URL.Path)},
		"npa": {"1"}, // Disabling Advertising Personalization

	}
	// event specific k/v
	for k, v := range *values {
		vv := make([]string, 1)
		vv[0] = url.QueryEscape(v)
		formValues[k] = vv
	}

	resp, err := http.PostForm(analyticsEndpoint, formValues)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Google Analytics returned '%d'", resp.StatusCode)
	}
	return nil
}
