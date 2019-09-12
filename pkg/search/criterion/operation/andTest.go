package operation

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	"github.com/BRBussy/bizzle/pkg/search/criterion/number"
	"github.com/stretchr/testify/suite"
)

type andTest struct {
	suite.Suite
}

func (t *andTest) Test() {
	testCriteria := And{}

	t.Equal(testCriteria.Type(), criterion.OperationAndCriterionType)

	t.Equal(testCriteria.IsValid(), criterion.ErrInvalid{Reasons: []string{"no criteria to and together"}})

	numberCriteria := number.Exact{
		Field:  "amountDue",
		Number: 1234,
	}

	testCriteria.Criteria = []criterion.Criterion{
		numberCriteria,
	}

	t.Equal(testCriteria.IsValid(), nil)

	t.Equal(
		testCriteria.ToFilter(),
		map[string]interface{}{
			"$and": []map[string]interface{}{
				numberCriteria.ToFilter(),
			},
		},
	)
}
