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

var fringeTestCases = []compareTestCase{
	{
		id:     "fringeTestCase 1",
		a:      nil,
		b:      nil,
		result: true,
	},
}

var fieldOnlyTestCases = []compareTestCase{
	{
		id:     "fieldOnlyTestCase 1",
		a:      make([]searchCriterion.Criterion, 0),
		b:      make([]searchCriterion.Criterion, 0),
		result: true,
	},
	{
		id: "fieldOnlyTestCase 2",
		a:  make([]searchCriterion.Criterion, 0),
		b: []searchCriterion.Criterion{
			stringCriterion.Substring{},
		},
		result: false,
	},
	{
		id: "fieldOnlyTestCase 3",
		a: []searchCriterion.Criterion{
			stringCriterion.Substring{},
		},
		b:      make([]searchCriterion.Criterion, 0),
		result: false,
	},
	{
		id: "fieldOnlyTestCase 4",
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
		id: "fieldOnlyTestCase 5",
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
		id: "fieldOnlyTestCase 6",
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
		id: "fieldOnlyTestCase 7",
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

var operationORTestCases = []compareTestCase{
	{
		id: "operationORTestCase 1",
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
		id: "operationORTestCase 2",
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
		id: "operationORTestCase 3",
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
		id: "operationORTestCase 4",
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
		result: false,
	},
	{
		id: "operationORTestCase 5",
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

var operationANDTestCases = []compareTestCase{
	{
		id: "operationANDTestCase 1",
		a: []searchCriterion.Criterion{
			operationCriterion.And{
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
			operationCriterion.And{
				Criteria: Criteria{},
			},
		},
		result: true,
	},
	{
		id: "operationANDTestCase 2",
		a: []searchCriterion.Criterion{
			operationCriterion.And{
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
			operationCriterion.And{
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
		id: "operationANDTestCase 3",
		a: []searchCriterion.Criterion{
			operationCriterion.And{
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
			operationCriterion.And{
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
		id: "operationANDTestCase 4",
		a: []searchCriterion.Criterion{
			operationCriterion.And{
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
			operationCriterion.And{
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
		result: false,
	},
	{
		id: "operationANDTestCase 5",
		a: []searchCriterion.Criterion{
			stringCriterion.Substring{
				Field:  "testStringSubstringField",
				String: "testStringSubstring",
			},
			operationCriterion.And{
				Criteria: Criteria{
					stringCriterion.Exact{
						Field:  "testStringExactField",
						String: "testStringExact",
					},
					operationCriterion.And{
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
			operationCriterion.And{
				Criteria: Criteria{
					operationCriterion.And{
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

var combinedTestCases = []compareTestCase{
	{
		id: "combinedTestCase 1",
		a: []searchCriterion.Criterion{
			stringCriterion.Substring{
				Field:  "testStringSubstringField",
				String: "testStringSubstring",
			},
			operationCriterion.Or{
				Criteria: Criteria{
					operationCriterion.And{
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
					operationCriterion.Or{
						Criteria: Criteria{
							operationCriterion.And{
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
			stringCriterion.Substring{
				Field:  "testStringExactField",
				String: "testStringExact",
			},
		},
		b: []searchCriterion.Criterion{
			operationCriterion.Or{
				Criteria: Criteria{
					operationCriterion.And{
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
					operationCriterion.Or{
						Criteria: Criteria{
							operationCriterion.And{
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
			stringCriterion.Substring{
				Field:  "testStringExactField",
				String: "testStringExact",
			},
			stringCriterion.Substring{
				Field:  "testStringSubstringField",
				String: "testStringSubstring",
			},
		},
		result: true,
	},
}

var allCompareTestCases = [][]compareTestCase{
	fringeTestCases,
	fieldOnlyTestCases,
	operationORTestCases,
	operationANDTestCases,
	combinedTestCases,
}

func TestCriteriaCompare(t *testing.T) {
	assert := testifyAssert.New(t)

	for i := range allCompareTestCases {
		for j := range allCompareTestCases[i] {
			assert.Equal(
				allCompareTestCases[i][j].result,
				Compare(allCompareTestCases[i][j].a, allCompareTestCases[i][j].b),
				fmt.Sprintf("%s", allCompareTestCases[i][j].id),
			)
		}
	}
}
