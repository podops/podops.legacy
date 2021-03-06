package platform

import (
	p "github.com/fupas/platform"
)

func ReportError(e error) {
	p.ReportError(e)
}

/*
var (
	// a central error client instance used in the service
	errorClient *o.Client
)

func init() {
	projectID := env.GetString("PROJECT_ID", "")
	if projectID == "" {
		log.Fatal("Missing variable 'PROJECT_ID'")
	}

	cl, err := o.NewClient(context.Background(), projectID, env.GetString("SERVICE_NAME", "default"))
	if err != nil {
		log.Fatal(err)
	}
	errorClient = cl
}

// ReportError reports an error, what else?
func ReportError(err error) {
	if errorClient == nil {
		log.Fatal(fmt.Errorf("Google Error Reporting is not initialized"))
	}
	errorClient.ErrorClient.Report(errorreporting.Entry{Error: err})
}
*/
