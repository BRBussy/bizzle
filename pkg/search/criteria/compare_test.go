package criteria

import (
	"fmt"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	numberCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/number"
	operationCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/operation"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/string"
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

type compareTestCase struct {
	id     string
	a      []searchCriterion.Criterion
	b      []searchCriterion.Criterion
	result bool
}

var compareTestCases = []compareTestCase{
	{
		id:     "0",
		a:      make([]searchCriterion.Criterion, 0),
		b:      make([]searchCriterion.Criterion, 0),
		result: true,
	},
	{
		id: "2",
		a:  make([]searchCriterion.Criterion, 0),
		b: []searchCriterion.Criterion{
			stringCriterion.Substring{},
		},
		result: false,
	},
	{
		id: "3",
		a: []searchCriterion.Criterion{
			stringCriterion.Substring{},
		},
		b:      make([]searchCriterion.Criterion, 0),
		result: false,
	},
	{
		id: "4",
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
		id: "5",
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
		id: "6",
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
		id: "7",
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
	{
		id: "8",
		a: []searchCriterion.Criterion{
			operationCriterion.Or{
				Criteria: Criteria{},
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
			operationCriterion.Or{
				Criteria: Criteria{},
			},
		},
		result: true,
	},
	{
		id: "9",
		a: []searchCriterion.Criterion{
			operationCriterion.Or{
				Criteria: Criteria{
					numberCriterion.Range{
						Field: "testNumberRange",
						Start: numberCriterion.RangeValue{
							Number:    123,
							Inclusive: true,
							Ignore:    false,
						},
						End: numberCriterion.RangeValue{
							Number:    100,
							Inclusive: false,
							Ignore:    true,
						},
					},
				},
			},
			stringCriterion.Substring{
				Field:  "testStringExactField",
				String: "testStringExact",
			},
		},
		b: []searchCriterion.Criterion{
			stringCriterion.Substring{
				Field:  "testStringExactField",
				String: "testStringExact",
			},
			operationCriterion.Or{
				Criteria: Criteria{
					stringCriterion.Substring{
						Field:  "testStringExactField",
						String: "testStringExact",
					},
				},
			},
		},
		result: false,
	},
	{
		id: "10",
		a: []searchCriterion.Criterion{
			operationCriterion.Or{
				Criteria: Criteria{
					numberCriterion.Range{
						Field: "testNumberRange",
						Start: numberCriterion.RangeValue{
							Number:    123,
							Inclusive: true,
							Ignore:    false,
						},
						End: numberCriterion.RangeValue{
							Number:    100,
							Inclusive: false,
							Ignore:    true,
						},
					},
					stringCriterion.Exact{
						Field:  "testStringExactField",
						String: "testStringExact",
					},
				},
			},
			stringCriterion.Substring{
				Field:  "testStringSubstringField",
				String: "testStringSubstring",
			},
		},
		b: []searchCriterion.Criterion{
			stringCriterion.Substring{
				Field:  "testStringSubstringField",
				String: "testStringSubstring",
			},
			operationCriterion.Or{
				Criteria: Criteria{
					stringCriterion.Exact{
						Field:  "testStringExactField",
						String: "testStringExact",
					},
					numberCriterion.Range{
						Field: "testNumberRange",
						Start: numberCriterion.RangeValue{
							Number:    123,
							Inclusive: true,
							Ignore:    false,
						},
						End: numberCriterion.RangeValue{
							Number:    100,
							Inclusive: false,
							Ignore:    true,
						},
					},
				},
			},
		},
		result: true,
	},
	{
		id: "11",
		a: []searchCriterion.Criterion{
			operationCriterion.Or{
				Criteria: Criteria{
					stringCriterion.Exact{
						Field:  "testStringExactField",
						String: "testStringExact",
					},
				},
			},
			stringCriterion.Substring{
				Field:  "testStringSubstringField",
				String: "testStringSubstring",
			},
		},
		b: []searchCriterion.Criterion{
			stringCriterion.Substring{
				Field:  "testStringSubstringField",
				String: "testStringSubstring",
			},
			operationCriterion.Or{
				Criteria: Criteria{
					stringCriterion.Exact{
						Field:  "testStringExactField",
						String: "testStringExact",
					},
					numberCriterion.Range{
						Field: "testNumberRange",
						Start: numberCriterion.RangeValue{
							Number:    123,
							Inclusive: true,
							Ignore:    false,
						},
						End: numberCriterion.RangeValue{
							Number:    100,
							Inclusive: false,
							Ignore:    true,
						},
					},
				},
			},
		},
		result: true,
	},
	{
		id: "12",
		a: []searchCriterion.Criterion{
			stringCriterion.Substring{
				Field:  "testStringSubstringField",
				String: "testStringSubstring",
			},
			operationCriterion.Or{
				Criteria: Criteria{
					stringCriterion.Exact{
						Field:  "testStringExactField",
						String: "testStringExact",
					},
					operationCriterion.Or{
						Criteria: Criteria{
							numberCriterion.Range{
								Field: "testNumberRange",
								Start: numberCriterion.RangeValue{
									Number:    123,
									Inclusive: true,
									Ignore:    false,
								},
								End: numberCriterion.RangeValue{
									Number:    100,
									Inclusive: false,
									Ignore:    true,
								},
							},
							stringCriterion.Exact{
								Field:  "testStringExactField",
								String: "testStringExact",
							},
						},
					},
					numberCriterion.Range{
						Field: "testNumberRange",
						Start: numberCriterion.RangeValue{
							Number:    123,
							Inclusive: true,
							Ignore:    false,
						},
						End: numberCriterion.RangeValue{
							Number:    100,
							Inclusive: false,
							Ignore:    true,
						},
					},
				},
			},
		},
		b: []searchCriterion.Criterion{
			stringCriterion.Substring{
				Field:  "testStringSubstringField",
				String: "testStringSubstring",
			},
			operationCriterion.Or{
				Criteria: Criteria{
					operationCriterion.Or{
						Criteria: Criteria{
							stringCriterion.Exact{
								Field:  "testStringExactField",
								String: "testStringExact",
							},
							numberCriterion.Range{
								Field: "testNumberRange",
								Start: numberCriterion.RangeValue{
									Number:    123,
									Inclusive: true,
									Ignore:    false,
								},
								End: numberCriterion.RangeValue{
									Number:    100,
									Inclusive: false,
									Ignore:    true,
								},
							},
						},
					},
					stringCriterion.Exact{
						Field:  "testStringExactField",
						String: "testStringExact",
					},
					numberCriterion.Range{
						Field: "testNumberRange",
						Start: numberCriterion.RangeValue{
							Number:    123,
							Inclusive: true,
							Ignore:    false,
						},
						End: numberCriterion.RangeValue{
							Number:    100,
							Inclusive: false,
							Ignore:    true,
						},
					},
				},
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
			fmt.Sprintf("case %s", compareTestCases[i].id),
		)
	}
}
