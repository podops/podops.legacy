package apiv1

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"

	"gopkg.in/yaml.v2"
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
	err = r.Validate()
	if err != nil {
		return nil, "", fmt.Errorf("Resource is not valid. Reason: %w", err)
	}

	return &r, r.GUID(), nil
}

func loadEpisodeResource(data []byte) (interface{}, string, error) {
	var r Episode

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
