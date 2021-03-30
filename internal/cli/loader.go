package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"

	"gopkg.in/yaml.v3"

	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/validator"
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
	resourceLoaders[a.ResourceShow] = loadShowResource
	resourceLoaders[a.ResourceEpisode] = loadEpisodeResource
}

// LoadResource takes a byte array and determines its kind before unmarshalling it into its struct form
func LoadResource(data []byte) (interface{}, string, string, error) {

	r, err := LoadResourceMetadata(data)
	loader := resourceLoaders[r.Kind]
	if loader == nil {
		return nil, "", "", fmt.Errorf("unsupported resource '%s'", r.Kind)
	}

	resource, guid, err := loader(data)
	if err != nil {
		return nil, "", "", err
	}
	return resource, r.Kind, guid, nil
}

// LoadResourceMetadata reads only the metadata of a resource
func LoadResourceMetadata(data []byte) (*a.ResourceMetadata, error) {
	var r a.ResourceMetadata

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, fmt.Errorf("can not parse resource: %w", err)
	}
	return &r, nil
}

func loadShowResource(data []byte) (interface{}, string, error) {
	var r a.Show

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, "", fmt.Errorf("can not parse resource: %w", err)
	}
	v := r.Validate(validator.New(a.ResourceShow))
	if !v.IsValid() {
		return nil, "", fmt.Errorf(v.Error())
	}

	return &r, r.GUID(), nil
}

func loadEpisodeResource(data []byte) (interface{}, string, error) {
	var r a.Episode

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, "", fmt.Errorf("can not parse resource: %w", err)
	}
	v := r.Validate(validator.New(a.ResourceEpisode))
	if !v.IsValid() {
		return nil, "", fmt.Errorf(v.Error())
	}

	return &r, r.GUID(), nil
}
