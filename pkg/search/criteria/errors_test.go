package criteria

import (
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestErrors(t *testing.T) {
	assert := testifyAssert.New(t)
	assert.Equal(
		ErrInvalidSerializedCriteria{Reasons: []string{"r1", "r2"}}.Error(),
		"serialized criteria is invalid: r1, r2",
	)
	assert.Equal(
		ErrUnmarshal{Reasons: []string{"r1", "r2"}}.Error(),
		"unmarshalling error: r1, r2",
	)
}
