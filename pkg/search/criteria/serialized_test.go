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
		"{\"someField\":{\"type\":\"%s\",\"string\":\"someSubstring\"}}",
		searchCriterion.StringSubstringCriterionType,
	)),
	criterion: stringCriterion.Substring{
		Field:  "someField",
		String: "someSubstring",
	},
}

var stringExact1 = testCase{
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"someField\":{\"type\":\"%s\",\"string\":\"someExactString\"}}",
		searchCriterion.StringExactCriterionType,
	)),
	criterion: stringCriterion.Exact{
		Field:  "someField",
		String: "someExactString",
	},
}

var numberRange1 = testCase{
	serializedCriterion: []byte(fmt.Sprintf(
		"{\"someField\":{\"type\":\"%s\",\"start\":{\"number\":123.12,\"inclusive\":true,\"ignore\":false},\"end\":{\"number\":245.123,\"inclusive\":false,\"ignore\":false}}}",
		searchCriterion.NumberRangeCriterionType,
	)),
	criterion: numberCriterion.Range{
		Field: "someField",
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
		"{\"someField\":{\"type\":\"%s\",\"number\":123.45}}",
		searchCriterion.NumberExactCriterionType,
	)),
	criterion: numberCriterion.Exact{
		Field:  "someField",
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
			"{\"$or\":[{\"someField\":{\"type\":\"%s\",\"string\":555},\"someOtherField\":{\"type\":\"%s\",\"string\":555}}]}",
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

func TestSerializedCriteriaOrCriterion(t *testing.T) {
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
	assert.Equal(
		[]searchCriterion.Criterion{
			operationCriterion.Or{
				Criteria: []searchCriterion.Criterion{
					stringSubstring1.criterion,
					stringExact1.criterion,
					and1.criterion,
				},
			},
		},
		testSerializedCriteria.Criteria,
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
	assert.Equal(
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
	)
}