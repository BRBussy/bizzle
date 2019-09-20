package identifier

import (
	"fmt"
	testtifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailIdentifier(t *testing.T) {
	assert := testtifyAssert.New(t)

	testIdentifier := Email("")

	// test type
	assert.Equal(
		EmailIdentifierType,
		testIdentifier.Type(),
		"type should be correct",
	)

	// test validity
	assert.EqualError(
		testIdentifier.IsValid(),
		ErrInvalidIdentifier{Reasons: []string{"Email identifier is blank"}}.Error(),
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
			"email": "Bob",
		},
		testIdentifier.ToFilter(),
		"to filter should return correct value",
	)

	// test MarshalJSON
	marshalledID, err := testIdentifier.ToJSON()
	assert.Equal(
		nil,
		err,
		"error should be nil calling to JSON on identifier",
	)
	assert.JSONEq(
		fmt.Sprintf(
			"{\"type\":\"%s\",\"email\":\"Bob\"}",
			EmailIdentifierType,
		),
		string(marshalledID),
		"json data from marshalled name should be correct",
	)
}
