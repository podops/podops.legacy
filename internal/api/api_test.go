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
