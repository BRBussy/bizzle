package criteria

import (
	"fmt"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	numberCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/number"
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
	{
		a: []searchCriterion.Criterion{
			stringCriterion.Substring{
				Field:  "testField",
				String: "testString",
			},
		},
		b: []searchCriterion.Criterion{
			stringCriterion.Substring{
				Field:  "testField",
				String: "testString",
			},
		},
		result: true,
	},
	{
		a: []searchCriterion.Criterion{
			stringCriterion.Substring{
				Field:  "testFieldDifferent",
				String: "testString",
			},
		},
		b: []searchCriterion.Criterion{
			stringCriterion.Substring{
				Field:  "testField",
				String: "testString",
			},
		},
		result: false,
	},
	{
		a: []searchCriterion.Criterion{
			stringCriterion.Substring{
				Field:  "testFieldDifferent",
				String: "testString",
			},
		},
		b: []searchCriterion.Criterion{
			stringCriterion.Substring{
				Field:  "testField",
				String: "testString",
			},
			stringCriterion.Substring{
				Field:  "testFieldDifferent",
				String: "testString",
			},
		},
		result: false,
	},
	{
		a: []searchCriterion.Criterion{
			numberCriterion.Exact{
				Field:  "testNumberExactField",
				Number: 112.123,
			},
			stringCriterion.Substring{
				Field:  "testSubstringField",
				String: "testSubstring",
			},
			stringCriterion.Exact{
				Field:  "testStringExactField",
				String: "testStringExact",
			},
		},
		b: []searchCriterion.Criterion{
			stringCriterion.Exact{
				Field:  "testStringExactField",
				String: "testStringExact",
			},
			numberCriterion.Exact{
				Field:  "testNumberExactField",
				Number: 112.123,
			},
			stringCriterion.Substring{
				Field:  "testSubstringField",
				String: "testSubstring",
			},
		},
		result: true,
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
