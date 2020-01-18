package text

import (
	"fmt"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestSubstring(t *testing.T) {
	assert := testifyAssert.New(t)

	testCriterion := Substring{}

	assert.Equal(
		criterion.TextSubstringCriterionType,
		testCriterion.Type(),
	)

	assert.Equal(
		criterion.ErrInvalid{Reasons: []string{
			"string is blank",
			"field is blank",
		}},
		testCriterion.IsValid(),
	)

	testCriterion.Text = "string"

	assert.Equal(
		criterion.ErrInvalid{Reasons: []string{
			"field is blank",
		}},
		testCriterion.IsValid(),
	)

	testCriterion.Text = ""
	testCriterion.Field = "someField"

	assert.Equal(
		criterion.ErrInvalid{Reasons: []string{
			"string is blank",
		}},
		testCriterion.IsValid(),
	)

	testCriterion.Text = "string"
	testCriterion.Field = "someField"

	assert.Equal(
		nil,
		testCriterion.IsValid(),
	)

	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$regex":   ".*string.*",
				"$options": "i",
			},
		},
		testCriterion.ToFilter(),
	)

	// confirm return value of ToJSON
	fieldName, jsonMessage, err := testCriterion.ToJSON()
	assert.Equal(nil, err)
	assert.Equal("someField", fieldName)
	assert.JSONEq(
		fmt.Sprintf(
			"{\"type\":\"%s\",\"string\":\"string\"}",
			criterion.TextSubstringCriterionType.String(),
		),
		string(jsonMessage),
	)
}
