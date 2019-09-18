package criteria

import (
	"fmt"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	numberCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/number"
	operationCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/operation"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/string"
)

//
// These are successful test cases. Failures are tested per case.
//

type serializedTestCase struct {
	id                  string
	serializedCriterion []byte
	criteria            []searchCriterion.Criterion
}

var emptyTestCase = serializedTestCase{
	id:                  "emptyTestCase",
	serializedCriterion: []byte("{}"),
	criteria:            make([]searchCriterion.Criterion, 0),
}

var stringSubstring1TestCase = serializedTestCase{
	id: "stringSubstringTestCase 1",
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"stringSubstring1Field\":{\"type\":\"%s\",\"string\":\"stringSubstring1TestCase\"}}",
		searchCriterion.StringSubstringCriterionType,
	)),
	criteria: []searchCriterion.Criterion{
		stringCriterion.Substring{
			Field:  "stringSubstring1Field",
			String: "stringSubstring1TestCase",
		},
	},
}

var stringExact1 = serializedTestCase{
	id: "stringExactTestCase 1",
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"stringExact1Field\":{\"type\":\"%s\",\"string\":\"stringExact1\"}}",
		searchCriterion.StringExactCriterionType,
	)),
	criteria: []searchCriterion.Criterion{
		stringCriterion.Exact{
			Field:  "stringExact1Field",
			String: "stringExact1",
		},
	},
}

var numberRange1 = serializedTestCase{
	id: "numberRangeTestCase 1",
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"numberRange1Field\":{\"type\":\"%s\",\"start\":{\"number\":123.12,\"inclusive\":true,\"ignore\":false},\"end\":{\"number\":245.123,\"inclusive\":false,\"ignore\":false}}}",
		searchCriterion.NumberRangeCriterionType,
	)),
	criteria: []searchCriterion.Criterion{
		numberCriterion.Range{
			Field: "numberRange1Field",
			Start: numberCriterion.RangeValue{
				Number:    123.12,
				Inclusive: true,
				Ignore:    false,
			},
			End: numberCriterion.RangeValue{
				Number:    245.123,
				Inclusive: false,
				Ignore:    false,
			},
		},
	},
}

var numberExact1 = serializedTestCase{
	id: "numberExactTestCase 1",
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"numberExact1Field\":{\"type\":\"%s\",\"number\":123.45}}",
		searchCriterion.NumberExactCriterionType,
	)),
	criteria: []searchCriterion.Criterion{
		numberCriterion.Exact{
			Field:  "numberExact1Field",
			Number: 123.45,
		},
	},
}

var operationAnd1 = serializedTestCase{
	id: "operationAndTestCase 1",
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"someField\":{\"type\":\"%s\",\"number\":123.45},\"someOtherField\":{\"type\":\"%s\",\"string\":\"someExactString\"}}",
		searchCriterion.NumberExactCriterionType,
		searchCriterion.StringExactCriterionType,
	)),
	criteria: []searchCriterion.Criterion{
		numberCriterion.Exact{
			Field:  "someField",
			Number: 123.45,
		},
		stringCriterion.Exact{
			Field:  "someOtherField",
			String: "someExactString",
		},
	},
}

var operationOrTestCase1 = serializedTestCase{
	id: "operationORTestCase 1",
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"$or\":[%s,%s,%s,%s,%s]}",
		stringSubstring1TestCase.serializedCriterion,
		stringExact1.serializedCriterion,
		numberRange1.serializedCriterion,
		numberExact1.serializedCriterion,
		fmt.Sprintf(
			"{\"someField\":{\"type\":\"%s\",\"number\":123.45},\"someOtherField\":{\"type\":\"%s\",\"string\":\"someExactString\"}}",
			searchCriterion.NumberExactCriterionType,
			searchCriterion.StringExactCriterionType,
		),
	)),
	criteria: []searchCriterion.Criterion{
		operationCriterion.Or{
			Criteria: []searchCriterion.Criterion{
				stringSubstring1TestCase.criteria[0],
				stringExact1.criteria[0],
				numberRange1.criteria[0],
				numberExact1.criteria[0],
				operationCriterion.And{
					Criteria: []searchCriterion.Criterion{
						numberCriterion.Exact{
							Field:  "someField",
							Number: 123.45,
						},
						stringCriterion.Exact{
							Field:  "someOtherField",
							String: "someExactString",
						},
					},
				},
			},
		},
	},
}

var operationOrTestCase2 = serializedTestCase{
	id: "operationORTestCase 2",
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"$or\":[%s,%s,%s]}",
		operationOrTestCase1.serializedCriterion,
		numberRange1.serializedCriterion,
		fmt.Sprintf(
			"{\"someField\":{\"type\":\"%s\",\"number\":123.45},\"someOtherField\":{\"type\":\"%s\",\"string\":\"someExactString\"}}",
			searchCriterion.NumberExactCriterionType,
			searchCriterion.StringExactCriterionType,
		),
	)),
	criteria: []searchCriterion.Criterion{
		operationCriterion.Or{
			Criteria: []searchCriterion.Criterion{
				operationOrTestCase1.criteria[0],
				numberRange1.criteria[0],
				operationCriterion.And{
					Criteria: []searchCriterion.Criterion{
						numberCriterion.Exact{
							Field:  "someField",
							Number: 123.45,
						},
						stringCriterion.Exact{
							Field:  "someOtherField",
							String: "someExactString",
						},
					},
				},
			},
		},
	},
}

var combinedTestCase1 = serializedTestCase{
	id: "combinedTestCase1",
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"$or\":[%s,%s,%s],\"stringSubstring1Field\":{\"type\":\"%s\",\"string\":\"stringSubstring1TestCase\"}}",
		operationOrTestCase1.serializedCriterion,
		numberRange1.serializedCriterion,
		fmt.Sprintf(
			"{\"someField\":{\"type\":\"%s\",\"number\":123.45},\"someOtherField\":{\"type\":\"%s\",\"string\":\"someExactString\"}}",
			searchCriterion.NumberExactCriterionType,
			searchCriterion.StringExactCriterionType,
		),
		searchCriterion.StringSubstringCriterionType,
	)),
	criteria: []searchCriterion.Criterion{
		operationCriterion.Or{
			Criteria: []searchCriterion.Criterion{
				operationOrTestCase1.criteria[0],
				numberRange1.criteria[0],
				operationCriterion.And{
					Criteria: []searchCriterion.Criterion{
						numberCriterion.Exact{
							Field:  "someField",
							Number: 123.45,
						},
						stringCriterion.Exact{
							Field:  "someOtherField",
							String: "someExactString",
						},
					},
				},
			},
		},
		stringCriterion.Substring{
			Field:  "stringSubstring1Field",
			String: "stringSubstring1TestCase",
		},
	},
}

var allMarshallUnmarshalTestCases = []serializedTestCase{
	emptyTestCase,
	stringSubstring1TestCase,
	stringExact1,
	numberRange1,
	numberExact1,
	operationAnd1,
	operationOrTestCase1,
	operationOrTestCase2,
	combinedTestCase1,
}
