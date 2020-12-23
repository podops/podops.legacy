package cli

// https://github.com/urfave/cli/blob/master/docs/v2/manual.md

import (
	"fmt"
	"io/ioutil"

	"github.com/podops/podops/pkg/metadata"
	"gopkg.in/yaml.v2"
)

func load(path string) ([]byte, error) {
	// FIXME: only local yaml is supported at the moment !

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Can not read file '%s'. %w", path, err)
	}
	return data, nil
}

func loadResource(path string) (interface{}, string, string, error) {
	data, err := load(path)
	if err != nil {
		return nil, "", "", fmt.Errorf("Can not read file '%s'. %w", path, err)
	}

	// peek into the resource to determin its type
	r, err := loadBasicResource(data)
	if err != nil {
		return nil, "", "", err
	}

	kind := r.(*metadata.BasicResource).Kind
	guid := r.(*metadata.BasicResource).Metadata.Labels[metadata.LabelGUID]

	loader := resourceLoaders[kind]
	if loader == nil {
		return nil, "", "", fmt.Errorf("Unsupported resource '%s'", kind)
	}
	resource, err := loader(data)

	if err != nil {
		return nil, "", "", err
	}

	return resource, kind, guid, nil
}

func loadBasicResource(data []byte) (interface{}, error) {
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
