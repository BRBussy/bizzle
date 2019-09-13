package string

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestExact(t *testing.T) {
	assert := testifyAssert.New(t)

	testCriterion := Exact{}

	assert.Equal(
		criterion.StringExactCriterionType,
		testCriterion.Type(),
	)

	assert.Equal(
		criterion.ErrInvalid{Reasons: []string{
			"string is blank",
			"field is blank",
		}},
		testCriterion.IsValid(),
	)

	testCriterion.String = "string"

	assert.Equal(
		criterion.ErrInvalid{Reasons: []string{
			"field is blank",
		}},
		testCriterion.IsValid(),
	)

	testCriterion.String = ""
	testCriterion.Field = "someField"

	assert.Equal(
		criterion.ErrInvalid{Reasons: []string{
			"string is blank",
		}},
		testCriterion.IsValid(),
	)

	testCriterion.String = "string"
	testCriterion.Field = "someField"

	assert.Equal(
		nil,
		testCriterion.IsValid(),
	)

	assert.Equal(
		map[string]interface{}{
			"someField": "string",
		},
		testCriterion.ToFilter(),
	)
}
