package podops

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleConfiguration(t *testing.T) {
	client, err := NewClient(context.TODO(), "po-xxx-xxx")

	if assert.NoError(t, err) {
		assert.NotNil(t, client)
	}
}
