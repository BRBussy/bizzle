package criteria

import (
	"fmt"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalSerializedCriteria_InvalidInput(t *testing.T) {
	assert := testifyAssert.New(t)

	// test fringe invalid json inputs
	assert.Error(
		(&Serialized{}).UnmarshalJSON(nil),
		ErrInvalidSerializedCriteria{Reasons: []string{"json criterion data is nil"}}.Error(),
	)
	assert.Error(
		(&Serialized{}).UnmarshalJSON([]byte("notValidJSON")),
		ErrUnmarshal{Reasons: []string{"json unmarshal", "invalid character 'o' in literal null (expecting 'u')"}}.Error(),
	)
}

func TestBlankInput(t *testing.T) {
	assert := testifyAssert.New(t)

	testSerializedCriteria := Serialized{}

	assert.Equal(
		nil,
		(&testSerializedCriteria).UnmarshalJSON([]byte("{\"serialized\":null,\"criteria\":null}")),
	)

	assert.Equal(
		make([]searchCriterion.Criterion, 0),
		testSerializedCriteria.Criteria,
	)
}

func TestUnmarshalSerializedCriteria_OROperatorFailures(t *testing.T) {
	assert := testifyAssert.New(t)
	// invalid value provided for $or operator
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"or array unmarshal",
			"json: cannot unmarshal object into Go value of type []map[string]json.RawMessage",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"$or\":{}}")),
	)
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"or array unmarshal",
			"json: cannot unmarshal number into Go value of type []map[string]json.RawMessage",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"$or\":6}")),
	)
	assert.Equal(
		ErrUnmarshal{Reasons: []string{
			"or array unmarshal",
			"json: cannot unmarshal string into Go value of type []map[string]json.RawMessage",
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

func TestUnmarshalSerializedCriteria_FieldCriterionFailures(t *testing.T) {
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

func TestUnmarshalSerializedCriteria_StringSubstringFailures(t *testing.T) {
	assert := testifyAssert.New(t)

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
}

func TestUnmarshalSerializedCriteria_StringExactFailures(t *testing.T) {
	assert := testifyAssert.New(t)
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
}

func TestUnmarshalSerializedCriteria_NumberRangeFailures(t *testing.T) {
	assert := testifyAssert.New(t)
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
}

func TestUnmarshalSerializedCriteria_NumberExactFailures(t *testing.T) {
	assert := testifyAssert.New(t)
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
}

func TestUnmarshalSerializedCriteria_UnmarshalSuccesses(t *testing.T) {
	assert := testifyAssert.New(t)
	for _, testCase := range allMarshallUnmarshalTestCases {
		testSerializedCriteria := Serialized{}
		assert.Equal(
			nil,
			(&testSerializedCriteria).UnmarshalJSON(testCase.serializedCriterion),
			fmt.Sprintf("%s: error should be nil after unmarshalling", testCase.id),
		)
		assert.Equal(
			true,
			Compare(
				testCase.criteria,
				testSerializedCriteria.Criteria,
			),
			fmt.Sprintf("%s: criteria should be correct", testCase.id),
		)
	}
}
