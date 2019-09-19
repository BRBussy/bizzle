package identifier

import (
	"fmt"
	testtifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestNameIdentifier(t *testing.T) {
	assert := testtifyAssert.New(t)

	testIdentifier := Name("")

	// test type
	assert.Equal(
		NameIdentifierType,
		testIdentifier.Type(),
		"type should be correct",
	)

	// test validity
	assert.EqualError(
		testIdentifier.IsValid(),
		ErrInvalidIdentifier{Reasons: []string{"Name identifier is blank"}}.Error(),
		"invalid reason should be correct",
	)

	// populate and test again
	testIdentifier = "Bob"
	assert.Equal(
		nil,
		testIdentifier.IsValid(),
	)

	// test filter
	assert.Equal(
		map[string]interface{}{
			"name": "Bob",
		},
		testIdentifier.ToFilter(),
		"to filter should return correct value",
	)

	// test MarshalJSON
	marshalledID, err := testIdentifier.MarshalJSON()
	assert.Equal(
		nil,
		err,
		"error should be nil calling to JSON on identifier",
	)
	assert.JSONEq(
		fmt.Sprintf(
			"{\"type\":\"%s\",\"name\":\"Bob\"}",
			NameIdentifierType,
		),
		string(marshalledID),
		"json data from marshalled name should be correct",
	)
}
