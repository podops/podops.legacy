package resources

import (
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"gopkg.in/yaml.v2"

	"github.com/txsvc/platform/pkg/platform"

	"github.com/podops/podops/internal/config"
	"github.com/podops/podops/pkg/metadata"
)

type (
	// ResourceLoaderFunc implements loading of resources
	ResourceLoaderFunc func(data []byte) (interface{}, string, error)
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

	bkt := platform.Storage().Bucket(config.BucketProduction)
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
	defer writer.Close()
	if _, err := writer.Write(data); err != nil {
		return err
	}

	return nil
}

// ReadResource reads a resource from Cloud Storage
func ReadResource(ctx context.Context, path string) (interface{}, string, string, error) {

	bkt := platform.Storage().Bucket(config.BucketProduction)
	reader, err := bkt.Object(path).NewReader(ctx)
	if err != nil {
		return nil, "", "", err
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, "", "", err
	}

	return LoadResource(data)
}

// LoadResource takes a byte array and determines its kind before unmarshalling it into its struct form
func LoadResource(data []byte) (interface{}, string, string, error) {

	r, err := LoadResourceMetadata(data)
	loader := resourceLoaders[r.Kind]
	if loader == nil {
		return nil, "", "", fmt.Errorf("Unsupported resource '%s'", r.Kind)
	}

	resource, guid, err := loader(data)
	if err != nil {
		return nil, "", "", err
	}
	return resource, r.Kind, guid, nil
}

// LoadResourceMetadata reads only the metadata of a resource
func LoadResourceMetadata(data []byte) (*metadata.ResourceMetadata, error) {
	var r metadata.ResourceMetadata

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, fmt.Errorf("Can not parse resource. %w", err)
	}
	return &r, nil
}

// Exists verifies the resource exists
func Exists(ctx context.Context, path string) bool {
	obj := platform.Storage().Bucket(config.BucketCDN).Object(path)
	_, err := obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false
	}
	return true
}

func loadShowResource(data []byte) (interface{}, string, error) {
	var r metadata.Show

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, "", fmt.Errorf("Can not parse resource. %w", err)
	}
	err = r.Validate()
	if err != nil {
		return nil, "", fmt.Errorf("Resource is not valid. Reason: %w", err)
	}

	return &r, r.GUID(), nil
}

func loadEpisodeResource(data []byte) (interface{}, string, error) {
	var r metadata.Episode

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, "", fmt.Errorf("Can not parse resource. %w", err)
	}
	err = r.Validate()
	if err != nil {
		return nil, "", fmt.Errorf("Resource is not valid. Reason: %w", err)
	}

	return &r, r.GUID(), nil
}
