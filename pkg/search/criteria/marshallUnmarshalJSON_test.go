package criteria

import (
	"fmt"
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalMarshallJSON(t *testing.T) {
	assert := testifyAssert.New(t)

	// test unmarshal -> marshall
	for _, testCase := range allMarshallUnmarshalTestCases {

		// first unmarshal
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

		// then marshall again
		serializedData, err := testSerializedCriteria.MarshalJSON()
		assert.Equal(
			nil,
			err,
			fmt.Sprintf("%s: error should be nil after marshalling", testCase.id),
		)
		assert.JSONEq(
			string(testCase.serializedCriterion),
			string(serializedData),
			fmt.Sprintf("%s: JSON value should be correct", testCase.id),
		)
	}
}

func TestMarshallUnmarshalJSON(t *testing.T) {
	assert := testifyAssert.New(t)

	// test marshall -> unmarshal
	for _, testCase := range allMarshallUnmarshalTestCases {

		// first marshall
		testSerializedCriteria := Serialized{
			Criteria: testCase.criteria,
		}
		serializedData, err := testSerializedCriteria.MarshalJSON()
		assert.Equal(
			nil,
			err,
			fmt.Sprintf("%s: error should be nil after marshalling", testCase.id),
		)
		assert.JSONEq(
			string(testCase.serializedCriterion),
			string(serializedData),
			fmt.Sprintf("%s: JSON value should be correct", testCase.id),
		)

		// then unmarshal
		assert.Equal(
			nil,
			(&testSerializedCriteria).UnmarshalJSON(serializedData),
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
