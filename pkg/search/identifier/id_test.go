package identifier

import (
	testtifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestIDIdentifier(t *testing.T) {
	assert := testtifyAssert.New(t)

	testIdentifier := ID("")

	// test type
	assert.Equal(
		IDIdentifierType,
		testIdentifier.Type(),
		"type should be correct",
	)

	// test validity
	assert.EqualError(
		testIdentifier.IsValid(),
		ErrInvalidIdentifier{Reasons: []string{"ID identifier is blank"}}.Error(),
		"invalid reason should be correct",
	)

	// populate and test again
	testIdentifier = "1234"
	assert.Equal(
		nil,
		testIdentifier.IsValid(),
	)

	// test filter
	assert.Equal(
		map[string]interface{}{
			"id": "1234",
		},
		testIdentifier.ToFilter(),
		"to filter should return correct value",
	)
}
