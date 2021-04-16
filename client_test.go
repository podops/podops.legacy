package podops

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	podcastName    string = "simple-podcast"
	podcastTitle   string = "Simple PodOps SDK Example"
	podcastSummary string = "A simple podcast for testing and experimentation. Created with the PodOps API."
)

func TestSimpleConfiguration(t *testing.T) {
	client, err := NewClient(context.TODO(), "po-xxx-xxx")

	if assert.NoError(t, err) {
		assert.NotNil(t, client)
	}
}

func TestSimpleLoadConfiguration(t *testing.T) {
	opts := LoadConfiguration()
	assert.NotNil(t, opts)

	client, err := NewClient(context.TODO(), opts.Token)
	if assert.NoError(t, err) {
		assert.NotNil(t, client)
	}
}
