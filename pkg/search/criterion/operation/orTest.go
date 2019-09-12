package operation

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	numberCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/number"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/string"
	"github.com/stretchr/testify/suite"
)

type orTest struct {
	suite.Suite
}

func (t *orTest) Test() {
	testCriteria := Or{}

	t.Equal(
		criterion.OperationOrCriterionType,
		testCriteria.Type(),
	)

	t.Equal(
		criterion.ErrInvalid{Reasons: []string{"no criteria to or together"}},
		testCriteria.IsValid(),
	)

	numberCriteria := numberCriterion.Exact{
		Field:  "amountDue",
		Number: 1234,
	}

	stringCriterion1 := stringCriterion.Exact{
		Field:  "name",
		String: "sam",
	}

	stringCriterion2 := stringCriterion.Substring{
		Field:  "surname",
		String: "smith",
	}

	testCriteria.Criteria = []criterion.Criterion{
		numberCriteria,
		Or{
			Criteria: []criterion.Criterion{
				stringCriterion1,
				stringCriterion2,
			},
		},
	}

	t.Equal(
		nil,
		testCriteria.IsValid(),
	)

	t.Equal(
		map[string]interface{}{
			"$or": []map[string]interface{}{
				numberCriteria.ToFilter(),
				{
					"$or": []map[string]interface{}{
						stringCriterion1.ToFilter(),
						stringCriterion2.ToFilter(),
					},
				},
			},
		},
		testCriteria.ToFilter(),
	)
}
