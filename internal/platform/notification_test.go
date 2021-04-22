package platform

import (
	"testing"

	"github.com/fupas/commons/pkg/env"
	"github.com/stretchr/testify/assert"
)

/*
must export these ENV variables for the tests to work:

export EMAIL_DOMAIN=
export EMAIL_API_KEY=
export EMAIL_FROM_AND_TO=

*/

func TestSendEmail(t *testing.T) {
	fromAndTo := env.GetString("EMAIL_FROM_AND_TO", "")

	err := SendEmail(fromAndTo, fromAndTo, "unit test", "just a unit test")

	assert.NoError(t, err)
}
