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

	t.Equal(testCriterion.Type(), criterion.StringSubstringCriterionType)

	t.Equal(testCriterion.IsValid(), criterion.ErrInvalid{Reasons: []string{
		"string is blank",
		"field is blank",
	}})

	testCriterion.String = "string"

	t.Equal(testCriterion.IsValid(), criterion.ErrInvalid{Reasons: []string{
		"field is blank",
	}})

	testCriterion.String = ""
	testCriterion.Field = "someField"

	t.Equal(testCriterion.IsValid(), criterion.ErrInvalid{Reasons: []string{
		"string is blank",
	}})

	testCriterion.String = "string"
	testCriterion.Field = "someField"

	t.Equal(testCriterion.IsValid(), nil)

	t.Equal(
		testCriterion.ToFilter(),
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$regex":   ".*string.*",
				"$options": "i",
			},
		},
	)
}
