package criterion

import (
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestErrors(t *testing.T) {
	assert := testifyAssert.New(t)

	assert.EqualError(
		ErrInvalid{
			Reasons: []string{"r1", "r2"},
		},
		"criterion is invalid: r1, r2",
	)
}
