package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Scenario 1: new account, login, account confirmation, token swap
func TestValidateNotEmpty(t *testing.T) {

	a := "Aa"
	b := "bb"
	c := ""

	assert.False(t, validateNotEmpty())
	assert.False(t, validateNotEmpty(""))
	assert.True(t, validateNotEmpty(a, b))
	assert.False(t, validateNotEmpty(a, c, b))
}

func TestScope(t *testing.T) {

	scope1 := "production:read,production:write,production:build"

	assert.False(t, hasScope("", ""))
	assert.False(t, hasScope(scope1, ""))
	assert.False(t, hasScope("", scopeResourceRead))

	assert.True(t, hasScope(scope1, scopeProductionRead))
	assert.False(t, hasScope(scope1, scopeResourceRead))
}
