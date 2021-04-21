package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	//fn = "cover.png"
	fn = "./soundcheck.mp3"
)

func TestMetadata(t *testing.T) {

	meta, err := ExtractMetadataFromFile(fn)
	assert.NoError(t, err)
	assert.NotNil(t, meta)

}
