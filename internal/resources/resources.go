package resources

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v2"

	"cloud.google.com/go/storage"

	"github.com/txsvc/platform/pkg/platform"
)

// WriteResource creates a resource. An existing resource will be overwritten if force==true
func WriteResource(ctx context.Context, path string, create, force bool, rsrc interface{}) error {

	exists := true

	bkt := platform.Storage().Bucket(bucketProduction)
	obj := bkt.Object(path)

	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		exists = false
	}

	//fmt.Println(fmt.Sprintf("Exists: %v, Create: %v, Force: %v", exists, create, force))

	// some logic mangling here ...
	if create && exists && !force { // create on an existing resource
		return fmt.Errorf("resource: '%s' already exists", path)
	}
	if !exists && !create && !force { // update on a missing resource
		return fmt.Errorf("resource: '%s' does not exists", path)
	}

	data, err := yaml.Marshal(rsrc)
	if err != nil {
		return err
	}

	writer := obj.NewWriter(ctx)
	if _, err := writer.Write(data); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}
