package identifier

import (
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestTypes(t *testing.T) {
	assert := testifyAssert.New(t)

	assert.Equal(
		"testType",
		Type("testType").String(),
		"string value should be correct",
	)
}
