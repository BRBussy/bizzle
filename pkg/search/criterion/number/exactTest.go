package number

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	"github.com/stretchr/testify/suite"
)

type exactTest struct {
	suite.Suite
}

func (t *exactTest) Test() {
	// create a blank criterion
	testCriterion := Exact{}

	// confirm is valid fails with field is blank
	t.Equal(
		criterion.ErrInvalid{Reasons: []string{"field is blank"}},
		testCriterion.IsValid(),
	)

	// confirm that type returns correct type
	t.Equal(
		criterion.NumberExactCriterionType,
		testCriterion.Type(),
	)

	// populate field and value
	testCriterion.Field = "someField"
	testCriterion.Number = 123.123

	// confirm is valid does not fail
	t.Equal(
		nil,
		testCriterion.IsValid(),
	)

	// confirm return value of ToFilter
	t.Equal(
		map[string]interface{}{
			"someField": 123.123,
		},
		testCriterion.ToFilter(),
	)
}
