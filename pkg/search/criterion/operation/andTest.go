package operation

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	numberCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/number"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/string"
	"github.com/stretchr/testify/suite"
)

type andTest struct {
	suite.Suite
}

func (t *andTest) Test() {
	testCriteria := And{}

	t.Equal(
		criterion.OperationAndCriterionType,
		testCriteria.Type(),
	)

	t.Equal(
		criterion.ErrInvalid{Reasons: []string{"and operation criterion has an empty criterion array"}},
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
		And{
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
			"$and": []map[string]interface{}{
				numberCriteria.ToFilter(),
				{
					"$and": []map[string]interface{}{
						stringCriterion1.ToFilter(),
						stringCriterion2.ToFilter(),
					},
				},
			},
		},
		testCriteria.ToFilter(),
	)
}
