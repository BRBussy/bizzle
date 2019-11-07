package text

import (
	"fmt"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestList(t *testing.T) {
	assert := testifyAssert.New(t)

	testCriterion := List{}

	assert.Equal(
		criterion.TextListCriterionType,
		testCriterion.Type(),
	)

	assert.Equal(
		criterion.ErrInvalid{Reasons: []string{
			"field is blank",
			"list is empty",
		}},
		testCriterion.IsValid(),
	)

	testCriterion.List = []string{"thing1", "thing2"}

	assert.Equal(
		criterion.ErrInvalid{Reasons: []string{
			"field is blank",
		}},
		testCriterion.IsValid(),
	)

	testCriterion.List = nil
	testCriterion.Field = "someField"

	assert.Equal(
		criterion.ErrInvalid{Reasons: []string{
			"list is empty",
		}},
		testCriterion.IsValid(),
	)

	testCriterion.List = []string{"thing1", "thing2"}

	assert.Equal(
		nil,
		testCriterion.IsValid(),
	)

	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$in": []string{"thing1", "thing2"},
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
			"{\"type\":\"%s\",\"list\":[\"thing1\",\"thing2\"]}",
			criterion.TextListCriterionType.String(),
		),
		string(jsonMessage),
	)
}
