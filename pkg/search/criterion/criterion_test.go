package criterion

import (
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestCriterion(t *testing.T) {
	assert := testifyAssert.New(t)

	var testType Type = "testType"

	assert.Equal(
		"testType",
		testType.String(),
	)
}
