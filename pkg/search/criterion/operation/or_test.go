package operation

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	numberCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/number"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/string"
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestOr(t *testing.T) {
	assert := testifyAssert.New(t)

	testCriteria := Or{}

	assert.Equal(
		criterion.OperationOrCriterionType,
		testCriteria.Type(),
	)

	assert.Equal(
		criterion.ErrInvalid{Reasons: []string{"or operation criterion has an empty criterion array"}},
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

	assert.Equal(
		nil,
		testCriteria.IsValid(),
	)

	assert.Equal(
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
