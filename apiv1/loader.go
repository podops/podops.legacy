package apiv1

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"

	"gopkg.in/yaml.v3"
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
	resourceLoaders[ResourceShow] = loadShowResource
	resourceLoaders[ResourceEpisode] = loadEpisodeResource
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
func LoadResourceMetadata(data []byte) (*ResourceMetadata, error) {
	var r ResourceMetadata

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, fmt.Errorf("Can not parse resource. %w", err)
	}
	return &r, nil
}

func loadShowResource(data []byte) (interface{}, string, error) {
	var r Show

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, "", fmt.Errorf("Can not parse resource. %w", err)
	}
	v := r.Validate(NewValidator(ResourceShow))
	if !v.IsValid() {
		return nil, "", fmt.Errorf(v.Error())
	}

	return &r, r.GUID(), nil
}

func loadEpisodeResource(data []byte) (interface{}, string, error) {
	var r Episode

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, "", fmt.Errorf("Can not parse resource. %w", err)
	}
	v := r.Validate(NewValidator(ResourceEpisode))
	if !v.IsValid() {
		return nil, "", fmt.Errorf(v.Error())
	}

	return &r, r.GUID(), nil
}
