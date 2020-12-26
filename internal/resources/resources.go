package resources

import (
	"context"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"cloud.google.com/go/storage"

	"github.com/podops/podops/pkg/metadata"
	"github.com/txsvc/platform/pkg/platform"
)

type (
	// ResourceLoaderFunc implements loading of resources
	ResourceLoaderFunc func(data []byte) (interface{}, error)
)

var (
	resourceLoaders map[string]ResourceLoaderFunc
)

func init() {

	resourceLoaders = make(map[string]ResourceLoaderFunc)
	resourceLoaders["show"] = loadShowResource
	resourceLoaders["episode"] = loadEpisodeResource
}

// WriteResource creates a resource. An existing resource will be overwritten if force==true
func WriteResource(ctx context.Context, path string, create, force bool, rsrc interface{}) error {

	exists := true

	bkt := platform.Storage().Bucket(bucketProduction)
	obj := bkt.Object(path)

	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		exists = false
	}

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

// ReadResource reads a resource
func ReadResource(ctx context.Context, path string) (interface{}, string, error) {

	bkt := platform.Storage().Bucket(bucketProduction)
	reader, err := bkt.Object(path).NewReader(ctx)
	if err != nil {
		return nil, "", err
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, "", err
	}
	r, err := loadBasicResource(data)
	loader := resourceLoaders[r.Kind]
	if loader == nil {
		return nil, "", fmt.Errorf("Unsupported resource '%s'", r.Kind)
	}

	resource, err := loader(data)
	if err != nil {
		return nil, "", err
	}
	return resource, r.Kind, nil
}

func loadBasicResource(data []byte) (*metadata.BasicResource, error) {
	var r metadata.BasicResource

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, fmt.Errorf("Can not parse resource. %w", err)
	}
	return &r, nil
}

func loadShowResource(data []byte) (interface{}, error) {
	var r metadata.Show

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, fmt.Errorf("Can not parse resource. %w", err)
	}
	err = r.Validate()
	if err != nil {
		return nil, fmt.Errorf("Resource is not valid. Reason: %w", err)
	}

	return &r, nil
}

func loadEpisodeResource(data []byte) (interface{}, error) {
	var r metadata.Episode

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, fmt.Errorf("Can not parse resource. %w", err)
	}
	err = r.Validate()
	if err != nil {
		return nil, fmt.Errorf("Resource is not valid. Reason: %w", err)
	}

	return &r, nil
}
