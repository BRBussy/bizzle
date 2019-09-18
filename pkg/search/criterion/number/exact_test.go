package number

import (
	"fmt"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestExact(t *testing.T) {
	assert := testifyAssert.New(t)

	// create a blank criterion
	testCriterion := Exact{}

	// confirm is valid fails with field is blank
	assert.Equal(
		criterion.ErrInvalid{Reasons: []string{"field is blank"}},
		testCriterion.IsValid(),
	)

	// confirm that type returns correct type
	assert.Equal(
		criterion.NumberExactCriterionType,
		testCriterion.Type(),
	)

	// populate field and value
	testCriterion.Field = "someField"
	testCriterion.Number = 123.123

	// confirm is valid does not fail
	assert.Equal(
		nil,
		testCriterion.IsValid(),
	)

	// confirm return value of ToFilter
	assert.Equal(
		map[string]interface{}{
			"someField": 123.123,
		},
		testCriterion.ToFilter(),
	)

	// confirm return value of ToJSON
	fieldName, jsonMessage, err := testCriterion.ToJSON()
	assert.Equal(nil, err)
	assert.Equal("someField", fieldName)
	assert.JSONEq(
		fmt.Sprintf(
			"{\"type\":\"%s\",\"number\":123.123}",
			criterion.NumberExactCriterionType.String(),
		),
		string(jsonMessage),
	)
}
