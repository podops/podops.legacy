package resources

import (
	"context"

	"github.com/txsvc/platform/pkg/platform"
	"gopkg.in/yaml.v2"
)

// CreateResource creates a resource. An existing resource will be overwritten if force==true
func CreateResource(ctx context.Context, path string, force bool, rsrc interface{}) error {
	// FIXME: implement force logic

	data, err := yaml.Marshal(rsrc)
	if err != nil {
		return err
	}

	bkt := platform.Storage().Bucket(bucketProduction)
	writer := bkt.Object(path).NewWriter(ctx)
	if _, err := writer.Write(data); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}
