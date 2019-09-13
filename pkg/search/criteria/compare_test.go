package criteria

import (
	"fmt"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/string"
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

type compareTestCase struct {
	a      []searchCriterion.Criterion
	b      []searchCriterion.Criterion
	result bool
}

var compareTestCases = []compareTestCase{
	{
		a:      make([]searchCriterion.Criterion, 0),
		b:      make([]searchCriterion.Criterion, 0),
		result: true,
	},
	{
		a: make([]searchCriterion.Criterion, 0),
		b: []searchCriterion.Criterion{
			stringCriterion.Substring{},
		},
		result: false,
	},
	{
		a: []searchCriterion.Criterion{
			stringCriterion.Substring{},
		},
		b:      make([]searchCriterion.Criterion, 0),
		result: false,
	},
}

func TestCriteriaCompare(t *testing.T) {
	assert := testifyAssert.New(t)

	for i := range compareTestCases {
		assert.Equal(
			compareTestCases[i].result,
			Compare(compareTestCases[i].a, compareTestCases[i].b),
			fmt.Sprintf("case %d", i),
		)
	}
}
