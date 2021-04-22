package metadata

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testFile     = "testfile.mp3"
	testFilePath = "./testfile.mp3"
	testFileExt  = ".mp3"
	parent       = "abcd"
)

func TestExtractMetadataFromFile(t *testing.T) {

	meta, err := ExtractMetadataFromFile(testFilePath)

	assert.NoError(t, err)
	assert.NotNil(t, meta)

	assert.Equal(t, meta.ContentType, "audio/mpeg")
	assert.Greater(t, meta.Duration, int64(0))
	assert.NotEmpty(t, meta.Etag, meta.Name, meta.Origin)
	assert.Equal(t, meta.Name, meta.Origin)
	assert.Greater(t, meta.Size, int64(0))
	assert.Greater(t, meta.Timestamp, int64(0))
	assert.Empty(t, meta.GUID, meta.ParentGUID)

	assert.True(t, meta.IsAudio())
}

func TestLocalNamePart(t *testing.T) {
	assert.Equal(t, LocalNamePart(testFilePath), testFile)
}

func TestFingerprintWithExt(t *testing.T) {
	fp := FingerprintWithExt(parent, testFilePath)

	assert.NotEmpty(t, fp)
	assert.True(t, strings.HasPrefix(fp, parent))
	assert.True(t, strings.HasSuffix(fp, testFileExt))
}
