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

type testCase struct {
	serializedCriterion []byte
	criterion           searchCriterion.Criterion
}

var stringSubstring1 = testCase{
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"stringSubstring1Field\":{\"type\":\"%s\",\"string\":\"stringSubstring1\"}}",
		searchCriterion.StringSubstringCriterionType,
	)),
	criterion: stringCriterion.Substring{
		Field:  "stringSubstring1Field",
		String: "stringSubstring1",
	},
}

var stringExact1 = testCase{
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"stringExact1Field\":{\"type\":\"%s\",\"string\":\"stringExact1\"}}",
		searchCriterion.StringExactCriterionType,
	)),
	criterion: stringCriterion.Exact{
		Field:  "stringExact1Field",
		String: "stringExact1",
	},
}

var numberRange1 = testCase{
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"numberRange1Field\":{\"type\":\"%s\",\"start\":{\"number\":123.12,\"inclusive\":true,\"ignore\":false},\"end\":{\"number\":245.123,\"inclusive\":false,\"ignore\":false}}}",
		searchCriterion.NumberRangeCriterionType,
	)),
	criterion: numberCriterion.Range{
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
}

var numberExact1 = testCase{
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"numberExact1Field\":{\"type\":\"%s\",\"number\":123.45}}",
		searchCriterion.NumberExactCriterionType,
	)),
	criterion: numberCriterion.Exact{
		Field:  "numberExact1Field",
		Number: 123.45,
	},
}

var and1 = testCase{
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"someField\":{\"type\":\"%s\",\"number\":123.45},\"someOtherField\":{\"type\":\"%s\",\"string\":\"someExactString\"}}",
		searchCriterion.NumberExactCriterionType,
		searchCriterion.StringExactCriterionType,
	)),
	criterion: operationCriterion.And{
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
}

func TestSerializedCriteriaInvalidInput(t *testing.T) {
	assert := testifyAssert.New(t)

	// test fringe invalid json inputs
	assert.Equal(
		ErrInvalidSerializedCriteria{Reasons: []string{"json criterion data is nil"}},
		(&Serialized{}).UnmarshalJSON(nil),
	)
	assert.Equal(
		ErrUnmarshal{Reasons: []string{"json unmarshal", "invalid character 'o' in literal null (expecting 'u')"}},
		(&Serialized{}).UnmarshalJSON([]byte("notValidJSON")),
	)

	// test empty
	testEmpty := Serialized{}
	assert.Equal(
		nil,
		(&testEmpty).UnmarshalJSON([]byte("{}")),
	)
	assert.Equal(
		make([]searchCriterion.Criterion, 0),
		testEmpty.Criteria,
	)
}

func TestSerializedCriteriaOROperatorFailures(t *testing.T) {
	assert := testifyAssert.New(t)
	// invalid value provided for $or operator
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"or array unmarshal",
			"json: cannot unmarshal object into Go value of type criteria.jsonObjectArray",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"$or\":{}}")),
	)
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"or array unmarshal",
			"json: cannot unmarshal number into Go value of type criteria.jsonObjectArray",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"$or\":6}")),
	)
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"or array unmarshal",
			"json: cannot unmarshal string into Go value of type criteria.jsonObjectArray",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"$or\":\"string\"}")),
	)
	// empty array
	assert.Equal(
		ErrInvalidSerializedCriteria{Reasons: []string{
			searchCriterion.ErrInvalid{Reasons: []string{"or operation criterion has an empty criterion array"}}.Error(),
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"$or\":[]}")),
	)
	// parsing failure in array
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"element in or",
			ErrUnmarshal{Reasons: []string{
				"string exact",
				"json: cannot unmarshal number into Go struct field Exact.string of type string",
			}}.Error(),
		}},
		(&Serialized{}).UnmarshalJSON([]byte(fmt.Sprintf(
			"{\"$or\":[{\"someField\":{\"type\":\"%s\",\"string\":555}}]}",
			searchCriterion.StringExactCriterionType,
		))),
	)
	// parsing failure in array with and
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"element in or",
			ErrUnmarshal{Reasons: []string{
				"string exact",
				"json: cannot unmarshal number into Go struct field Exact.string of type string",
			}}.Error(),
		}},
		(&Serialized{}).UnmarshalJSON([]byte(fmt.Sprintf(
			"{\"$or\":[{\"someField\":{\"type\":\"%s\",\"string\":555},\"someOtherField\":{\"type\":\"%s\",\"string\":\"someSubstring\"}}]}",
			searchCriterion.StringExactCriterionType,
			searchCriterion.StringSubstringCriterionType,
		))),
	)
}

func TestSerializedCriteriaFieldCriterionFailures(t *testing.T) {
	assert := testifyAssert.New(t)

	// invalid input for field criterion
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"criterion object unmarshal",
			"json: cannot unmarshal array into Go value of type criteria.typeHolder",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"someField\":[]}")),
	)
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"criterion object unmarshal",
			"json: cannot unmarshal number into Go value of type criteria.typeHolder",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"someField\":5}")),
	)
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"criterion object unmarshal",
			"json: cannot unmarshal string into Go value of type criteria.typeHolder",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"someField\":\"string\"}")),
	)
	// invalid type for field criterion
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"criterion object unmarshal",
			"json: cannot unmarshal number into Go struct field typeHolder.type of type criterion.Type",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"someField\":{\"type\":4}}")),
	)
	// invalid value for field criterion
	assert.Equal(
		ErrInvalidSerializedCriteria{Reasons: []string{
			"invalid field criterion type",
			"invalidType",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"someField\":{\"type\":\"invalidType\"}}")),
	)
}

func TestSerializedCriteriaStringSubstringCriterion(t *testing.T) {
	assert := testifyAssert.New(t)

	// unmarshalling failure
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"string substring",
			"json: cannot unmarshal number into Go struct field Substring.string of type string",
		}},
		(&Serialized{}).UnmarshalJSON([]byte(fmt.Sprintf(
			"{\"someField\":{\"type\":\"%s\",\"string\":555}}",
			searchCriterion.StringSubstringCriterionType,
		))),
	)

	// unmarshalling success
	testSerializedCriteria := Serialized{}
	assert.Equal(
		nil,
		(&testSerializedCriteria).UnmarshalJSON(stringSubstring1.serializedCriterion),
	)
	assert.Equal(
		[]searchCriterion.Criterion{
			stringSubstring1.criterion,
		},
		testSerializedCriteria.Criteria,
	)
}

func TestSerializedCriteriaStringExactCriterion(t *testing.T) {
	assert := testifyAssert.New(t)
	// unmarshalling failure
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"string exact",
			"json: cannot unmarshal number into Go struct field Exact.string of type string",
		}},
		(&Serialized{}).UnmarshalJSON([]byte(fmt.Sprintf(
			"{\"someField\":{\"type\":\"%s\",\"string\":555}}",
			searchCriterion.StringExactCriterionType,
		))),
	)

	// unmarshalling success
	testSerializedCriteria := Serialized{}
	assert.Equal(
		nil,
		(&testSerializedCriteria).UnmarshalJSON(stringExact1.serializedCriterion),
	)
	assert.Equal(
		[]searchCriterion.Criterion{
			stringExact1.criterion,
		},
		testSerializedCriteria.Criteria,
	)
}

func TestSerializedCriteriaNumberRangeCriterion(t *testing.T) {
	assert := testifyAssert.New(t)
	// unmarshalling failure
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"number range",
			"json: cannot unmarshal string into Go struct field RangeValue.start.number of type float64",
		}},
		(&Serialized{}).UnmarshalJSON([]byte(fmt.Sprintf(
			"{\"someField\":{\"type\":\"%s\",\"start\":{\"number\":\"123.12\",\"inclusive\":true,\"ignore\":false},\"end\":{\"number\":\"245.123\",\"inclusive\":false,\"ignore\":false}}}",
			searchCriterion.NumberRangeCriterionType,
		))),
	)

	// unmarshalling success
	testSerializedCriteria := Serialized{}
	assert.Equal(
		nil,
		(&testSerializedCriteria).UnmarshalJSON(numberRange1.serializedCriterion),
	)
	assert.Equal(
		[]searchCriterion.Criterion{
			numberRange1.criterion,
		},
		testSerializedCriteria.Criteria,
	)
}

func TestSerializedCriteriaNumberExactCriterion(t *testing.T) {
	assert := testifyAssert.New(t)
	// unmarshalling failure
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"number exact",
			"json: cannot unmarshal string into Go struct field Exact.number of type float64",
		}},
		(&Serialized{}).UnmarshalJSON([]byte(fmt.Sprintf(
			"{\"someField\":{\"type\":\"%s\",\"number\":\"123.45\"}}",
			searchCriterion.NumberExactCriterionType,
		))),
	)

	// unmarshalling success
	testSerializedCriteria := Serialized{}
	assert.Equal(
		nil,
		(&testSerializedCriteria).UnmarshalJSON(numberExact1.serializedCriterion),
	)
	assert.Equal(
		[]searchCriterion.Criterion{
			numberExact1.criterion,
		},
		testSerializedCriteria.Criteria,
	)
}

func TestSerializedCriteriaORCriterion(t *testing.T) {
	assert := testifyAssert.New(t)
	serializedValue := []byte(fmt.Sprintf(
		"{\"$or\":[%s,%s,%s]}",
		stringSubstring1.serializedCriterion,
		stringExact1.serializedCriterion,
		and1.serializedCriterion,
	))

	testSerializedCriteria := Serialized{}

	assert.Equal(
		nil,
		(&testSerializedCriteria).UnmarshalJSON(serializedValue),
	)

	expectedCriterion := []searchCriterion.Criterion{
		operationCriterion.Or{
			Criteria: []searchCriterion.Criterion{
				stringSubstring1.criterion,
				stringExact1.criterion,
				and1.criterion,
			},
		},
	}

	// NOTE: a direct equal assertion is not possible here since as the
	// order of keys in a map cannot be guaranteed the order of elements
	// in the resultant criteria slice may differ to the expected

	assert.Equal(
		len(expectedCriterion),
		len(testSerializedCriteria.Criteria),
		"criteria should have correct no. of entries",
	)
	assert.IsType(
		expectedCriterion[0],
		testSerializedCriteria.Criteria[0],
		"both resultant and expected first elements should have the same type",
	)
	// infer or type of expected and resultant first elements
	expectedORCriterion, ok := expectedCriterion[0].(operationCriterion.Or)
	assert.Equal(
		true,
		ok,
	)
	resultORCriterion, ok := testSerializedCriteria.Criteria[0].(operationCriterion.Or)
	assert.Equal(
		true,
		ok,
	)

	assert.Equal(
		len(expectedORCriterion.Criteria),
		len(resultORCriterion.Criteria),
		"the resultant OR criteria should have correct no of elements",
	)

	// check first and second element equality
	assert.Equal(
		expectedORCriterion.Criteria[0],
		resultORCriterion.Criteria[0],
	)
	assert.Equal(
		expectedORCriterion.Criteria[1],
		resultORCriterion.Criteria[1],
	)

	// confirm that 3rd elements have same type
	assert.IsType(
		expectedORCriterion.Criteria[2],
		resultORCriterion.Criteria[2],
	)
	// infer or type of expected and resultant 3rd elements
	expectedANDCriterion, ok := expectedORCriterion.Criteria[2].(operationCriterion.And)
	assert.Equal(
		true,
		ok,
	)
	resultANDCriterion, ok := resultORCriterion.Criteria[2].(operationCriterion.And)
	assert.Equal(
		true,
		ok,
	)
	assert.Equal(
		len(expectedANDCriterion.Criteria),
		len(resultANDCriterion.Criteria),
		"the resultant AND criteria should have correct no of elements",
	)
	// confirm contents of and are the same
	assert.ElementsMatch(
		expectedANDCriterion.Criteria,
		resultANDCriterion.Criteria,
	)
}

func TestSerializedCriteriaCombinedCriterion(t *testing.T) {
	assert := testifyAssert.New(t)
	serializedValue := []byte(fmt.Sprintf(
		"{\"$or\":[%s,%s,%s],\"someField\":{\"type\":\"%s\",\"string\":\"someSubstring\"}}",
		stringSubstring1.serializedCriterion,
		stringExact1.serializedCriterion,
		and1.serializedCriterion,
		searchCriterion.StringSubstringCriterionType,
	))

	testSerializedCriteria := Serialized{}

	assert.Equal(
		nil,
		(&testSerializedCriteria).UnmarshalJSON(serializedValue),
	)
	t.Logf("tthis test: %t", compareCriteria(
		[]searchCriterion.Criterion{
			operationCriterion.Or{
				Criteria: []searchCriterion.Criterion{
					stringSubstring1.criterion,
					stringExact1.criterion,
					and1.criterion,
				},
			},
			stringCriterion.Substring{
				Field:  "someField",
				String: "someSubstring",
			},
		},
		testSerializedCriteria.Criteria,
		assert,
	))
	assert.Equal(
		true,
		compareCriteria(
			[]searchCriterion.Criterion{
				operationCriterion.Or{
					Criteria: []searchCriterion.Criterion{
						stringSubstring1.criterion,
						stringExact1.criterion,
						and1.criterion,
					},
				},
				stringCriterion.Substring{
					Field:  "someField",
					String: "someSubstring",
				},
			},
			testSerializedCriteria.Criteria,
			assert,
		),
		"criteria should be equal",
	)
}

func compareCriteria(a, b []searchCriterion.Criterion, assert *testifyAssert.Assertions) bool {
	// check lengths
	if !assert.Equal(
		len(a),
		len(b),
		"lengths of given criteria differ",
	) {
		return false
	}

	// for every element in a
nextA:
	for ia := range a {
		// look through b for a match
		for ib := range b {
			// if a match is found go to next element ia
			switch typedA := a[ia].(type) {
			case operationCriterion.And:
				if compareANDCriterion(typedA, b[ib], assert) {
					continue nextA
				}
			case operationCriterion.Or:
				if compareORCriterion(typedA, b[ib], assert) {
					continue nextA
				}
			default:
				if assert.Equal(a[ia], b[ib]) {
					continue nextA
				}
			}
		}
		// if execution reaches here ia was not found in b
		return false
	}
	// if execution reaches here every ia was found in b
	return true
}

func compareANDCriterion(a operationCriterion.And, b searchCriterion.Criterion, assert *testifyAssert.Assertions) bool {
	// check that a and b are both and criterion
	typedB, ok := b.(operationCriterion.And)
	if !ok {
		return false
	}
	// check lengths of a and b are the same
	if !assert.Equal(
		len(a.Criteria),
		len(typedB.Criteria),
		"no of elements is not the same in a and b AND criteria",
	) {
		return false
	}
	return compareCriteria(a.Criteria, typedB.Criteria, assert)
}

func compareORCriterion(a operationCriterion.Or, b searchCriterion.Criterion, assert *testifyAssert.Assertions) bool {
	// check that a and b are both or criterion
	typedB, ok := b.(operationCriterion.Or)
	if !ok {
		return false
	}
	// check lengths of a and b are the same
	if !assert.Equal(
		len(a.Criteria),
		len(typedB.Criteria),
		"no of elements is not the same in a and b OR criteria",
	) {
		return false
	}
	return compareCriteria(a.Criteria, typedB.Criteria, assert)
}
