package string

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	"github.com/stretchr/testify/suite"
)

type substringTest struct {
	suite.Suite
}

func (t *substringTest) Test() {
	testCriterion := Substring{}

	t.Equal(
		criterion.StringSubstringCriterionType,
		testCriterion.Type(),
	)

	t.Equal(
		criterion.ErrInvalid{Reasons: []string{
			"string is blank",
			"field is blank",
		}},
		testCriterion.IsValid(),
	)

	testCriterion.String = "string"

	t.Equal(
		criterion.ErrInvalid{Reasons: []string{
			"field is blank",
		}},
		testCriterion.IsValid(),
	)

	testCriterion.String = ""
	testCriterion.Field = "someField"

	t.Equal(
		criterion.ErrInvalid{Reasons: []string{
			"string is blank",
		}},
		testCriterion.IsValid(),
	)

	testCriterion.String = "string"
	testCriterion.Field = "someField"

	t.Equal(
		nil,
		testCriterion.IsValid(),
	)

	t.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$regex":   ".*string.*",
				"$options": "i",
			},
		},
		testCriterion.ToFilter(),
	)
}
