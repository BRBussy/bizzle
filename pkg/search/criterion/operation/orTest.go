package operation

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	"github.com/BRBussy/bizzle/pkg/search/criterion/number"
	"github.com/stretchr/testify/suite"
)

type orTest struct {
	suite.Suite
}

func (t *orTest) Test() {
	testCriteria := Or{}

	t.Equal(testCriteria.Type(), criterion.OperationOrCriterionType)

	t.Equal(testCriteria.IsValid(), criterion.ErrInvalid{Reasons: []string{"no criteria to or together"}})

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
			"$or": []map[string]interface{}{
				numberCriteria.ToFilter(),
			},
		},
	)
}
