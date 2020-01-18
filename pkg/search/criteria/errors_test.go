package criteria

import (
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestErrors(t *testing.T) {
	assert := testifyAssert.New(t)
	assert.EqualError(
		ErrInvalidSerializedCriteria{Reasons: []string{"r1", "r2"}},
		"serialized criteria is invalid: r1, r2",
	)
	assert.EqualError(
		ErrUnmarshal{Reasons: []string{"r1", "r2"}},
		"unmarshalling error: r1, r2",
	)
}
