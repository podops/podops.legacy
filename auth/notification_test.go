package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txsvc/platform/pkg/env"
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
