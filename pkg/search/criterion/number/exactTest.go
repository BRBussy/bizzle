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
	t.Equal(testCriterion.IsValid(), criterion.ErrInvalid{Reasons: []string{"field is blank"}})

	// confirm that type returns correct type
	t.Equal(testCriterion.Type(), criterion.NumberExactCriterionType)

	// populate field and value
	testCriterion.Field = "someField"
	testCriterion.Number = 123.123

	// confirm is valid does not fail
	t.Equal(testCriterion.IsValid(), nil)

	// confirm return value of ToFilter
	t.Equal(
		testCriterion.ToFilter(),
		map[string]interface{}{
			"someField": 123.123,
		},
	)
}
