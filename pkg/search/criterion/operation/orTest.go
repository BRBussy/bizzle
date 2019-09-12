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

	t.Equal(testCriteria.Type(), criterion.OperationOrCriterionType)

	t.Equal(testCriteria.IsValid(), criterion.ErrInvalid{Reasons: []string{"no criteria to or together"}})

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

	t.Equal(testCriteria.IsValid(), nil)

	t.Equal(
		testCriteria.ToFilter(),
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
	)
}
