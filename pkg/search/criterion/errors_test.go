package criterion

import (
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestErrors(t *testing.T) {
	assert := testifyAssert.New(t)

	assert.Equal(
		"criterion is invalid: r1, r2",
		ErrInvalid{
			Reasons: []string{"r1", "r2"},
		}.Error(),
	)
}
