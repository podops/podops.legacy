package podops

import (
	"testing"
)

func TestNewClientFromFile(t *testing.T) {
	client, err := NewClientFromFile(DefaultConfigLocation())
	if err != nil {
		t.Error(err)
	}
	if err := client.Validate(); err != nil {
		t.Error(err)
	}
}
