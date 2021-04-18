package platform

import (
	"context"
	"encoding/json"
	"fmt"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"

	"github.com/fupas/commons/pkg/env"
)

// CreateTask is used to schedule a background task using the default queue.
// The payload can be any struct and will be marshalled into a json string.
func CreateTask(ctx context.Context, handler string, payload interface{}) (*taskspb.Task, error) {

	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		// observer.ReportError(err) FIXME this is just disabled
		return nil, err
	}
	defer client.Close()

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", env.GetString("PROJECT_ID", ""), env.GetString("LOCATION_ID", ""), env.GetString("DEFAULT_QUEUE", ""))

	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
					HttpMethod:  taskspb.HttpMethod_POST,
					RelativeUri: handler,
					Headers:     map[string]string{"Content-Type": "application/json"},
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
		req.Task.GetAppEngineHttpRequest().Body = b
	}

	task, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, err
	}

	return task, nil
}
