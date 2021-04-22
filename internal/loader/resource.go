package loader

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/podops/podops"
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
	resourceLoaders[podops.ResourceShow] = loadShowResource
	resourceLoaders[podops.ResourceEpisode] = loadEpisodeResource
}

// UnmarshalResource takes a byte array and determines its kind before unmarshalling it into its struct form
func UnmarshalResource(data []byte) (interface{}, string, string, error) {

	r, _ := LoadGenericResource(data)
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

// LoadGenericResource reads only the metadata of a resource
func LoadGenericResource(data []byte) (*podops.GenericResource, error) {
	var r podops.GenericResource

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		return nil, fmt.Errorf("can not parse resource: %w", err)
	}
	return &r, nil
}

func loadShowResource(data []byte) (interface{}, string, error) {
	var show podops.Show

	err := yaml.Unmarshal([]byte(data), &show)
	if err != nil {
		return nil, "", fmt.Errorf("can not parse resource: %w", err)
	}

	// FIXME validate before write, not on every read!

	/*
		v := r.Validate(validator.New(podops.ResourceShow))
		if !v.IsValid() {
			return nil, "", fmt.Errorf(v.Error())
		}
	*/
	return &show, show.GUID(), nil
}

func loadEpisodeResource(data []byte) (interface{}, string, error) {
	var episode podops.Episode

	err := yaml.Unmarshal([]byte(data), &episode)
	if err != nil {
		return nil, "", fmt.Errorf("can not parse resource: %w", err)
	}

	// FIXME validate before write, not on every read!

	/*
		v := r.Validate(validator.New(podops.ResourceEpisode))
		if !v.IsValid() {
			return nil, "", fmt.Errorf(v.Error())
		}
	*/

	return &episode, episode.GUID(), nil
}
