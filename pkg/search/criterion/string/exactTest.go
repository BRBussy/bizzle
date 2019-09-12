package string

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	"github.com/stretchr/testify/suite"
)

type exactTest struct {
	suite.Suite
}

func (t *exactTest) Test() {
	testCriterion := Exact{}

	t.Equal(
		criterion.StringExactCriterionType,
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
			"someField": "string",
		},
		testCriterion.ToFilter(),
	)
}
