package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	scopeProductionRead  = "production:read"
	scopeProductionWrite = "production:write"
	scopeProductionBuild = "production:build"
	scopeResourceRead    = "resource:read"
	scopeResourceWrite   = "resource:write"
)

func TestScope(t *testing.T) {

	scope1 := "production:read,production:write,production:build"

	assert.False(t, hasScope("", ""))
	assert.False(t, hasScope(scope1, ""))
	assert.False(t, hasScope("", scopeResourceRead))

	assert.True(t, hasScope(scope1, scopeProductionRead))
	assert.False(t, hasScope(scope1, scopeResourceRead))
}
