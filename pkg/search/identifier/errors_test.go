package identifier

import (
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestErrors(t *testing.T) {
	assert := testifyAssert.New(t)

	assert.EqualError(
		ErrInvalidIdentifier{Reasons: []string{"r1", "r2"}},
		"invalid identifier: r1, r2",
		"error message should be correct for error type",
	)
}
