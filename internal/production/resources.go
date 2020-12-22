package production

import (
	"context"

	"github.com/podops/podops/internal/errors"
	"github.com/txsvc/platform/pkg/platform"
	"gopkg.in/yaml.v2"
)

// CreateResource creates a resource. An existing resource will be overwritten if force==true
func CreateResource(ctx context.Context, path string, force bool, rsrc interface{}) error {
	// FIXME: implement force logic

	data, err := yaml.Marshal(rsrc)
	if err != nil {
		return errors.Wrap(err)
	}

	bkt := platform.Storage().Bucket(bucketProduction)
	writer := bkt.Object(path).NewWriter(ctx)
	if _, err := writer.Write(data); err != nil {
		return errors.Wrap(err)
	}
	if err := writer.Close(); err != nil {
		return errors.Wrap(err)
	}

	return nil
}
