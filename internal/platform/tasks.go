package platform

import (
	"context"
	"encoding/json"
	"fmt"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"

	"github.com/fupas/commons/pkg/env"
)

var (
	// UserAgentString identifies any http request podops makes
	userAgentString string = fmt.Sprintf("PodOps %d.%d.%d", 1, 0, 1) // FIXME this has to be aligned with version.go
	// workerQueue is the main worker queue for all the background tasks
	workerQueue string = fmt.Sprintf("projects/%s/locations/%s/queues/%s", env.GetString("PROJECT_ID", ""), env.GetString("LOCATION_ID", ""), env.GetString("DEFAULT_QUEUE", ""))
)

// CreateHttpTask is used to schedule a background task using the default queue.
// The payload can be any struct and will be marshalled into a json string.
func CreateHttpTask(ctx context.Context, method taskspb.HttpMethod, handler, token string, payload interface{}) (*taskspb.Task, error) {

	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		ReportError(err)
		return nil, err
	}
	defer client.Close()

	req := &taskspb.CreateTaskRequest{
		Parent: workerQueue,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: method,
					Url:        handler,
					Headers: map[string]string{
						"Content-Type":  "application/json",
						"User-Agent":    userAgentString,
						"Authorization": fmt.Sprintf("Bearer %s", token),
					},
				},
			},
		},
	}

	if payload != nil {
		// marshal the payload
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		req.Task.GetHttpRequest().Body = b
	}

	task, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, err
	}

	return task, nil
}
