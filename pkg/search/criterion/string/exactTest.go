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

	t.Equal(testCriterion.Type(), criterion.StringExactCriterionType)

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
}
