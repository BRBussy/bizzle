package criteria

import (
	"fmt"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	numberCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/number"
	operationCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/operation"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/text"
	testifyAssert "github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestSerializedCriteria_InvalidAndPlacement(t *testing.T) {
	assert := testifyAssert.New(t)

	// test incorrectly placed top level and
	testSerializedCriteria := Serialized{
		Criteria: []searchCriterion.Criterion{
			operationCriterion.And{},
		},
	}
	jsonData, err := testSerializedCriteria.MarshalJSON()
	assert.Equal(
		[]byte(nil),
		jsonData,
		"json data should be nil with marshalling error",
	)
	assert.EqualError(
		err,
		ErrInvalidSerializedCriteria{
			Reasons: []string{"and only allowed if contained within or"},
		}.Error(),
	)

	// test incorrectly placed embedded and
	testSerializedCriteria.Criteria = []searchCriterion.Criterion{
		operationCriterion.Or{
			Criteria: []searchCriterion.Criterion{
				operationCriterion.And{
					Criteria: []searchCriterion.Criterion{
						operationCriterion.And{},
					},
				},
			},
		},
	}
	jsonData, err = testSerializedCriteria.MarshalJSON()
	assert.Equal(
		[]byte(nil),
		jsonData,
		"json data should be nil with marshalling error",
	)
	assert.EqualError(
		err,
		ErrMarshal{
			Reasons: []string{
				"or element",
				ErrInvalidSerializedCriteria{
					Reasons: []string{"and only allowed if contained within or"},
				}.Error(),
			},
		}.Error(),
	)
}

func TestSerializedCriteria_MarshallFailures(t *testing.T) {
	assert := testifyAssert.New(t)

	// top level marshalling failure
	testSerializedCriteria := Serialized{
		Criteria: []searchCriterion.Criterion{
			numberCriterion.Exact{
				Field:  "testField",
				Number: math.Inf(1),
			},
		},
	}
	jsonData, err := testSerializedCriteria.MarshalJSON()
	assert.Equal(
		[]byte(nil),
		jsonData,
		"json data should be nil with marshalling error",
	)
	assert.EqualError(
		err,
		ErrMarshal{
			Reasons: []string{
				"json: unsupported value: +Inf",
			},
		}.Error(),
	)

	// embedded marshalling failure
	testSerializedCriteria.Criteria = []searchCriterion.Criterion{
		stringCriterion.Exact{
			Field:  "testField",
			String: "exactString",
		},
		operationCriterion.Or{
			Criteria: []searchCriterion.Criterion{
				numberCriterion.Exact{
					Field:  "testField",
					Number: math.Inf(1),
				},
			},
		},
	}
	jsonData, err = testSerializedCriteria.MarshalJSON()
	assert.Equal(
		[]byte(nil),
		jsonData,
		"json data should be nil with marshalling error",
	)
	assert.EqualError(
		err,
		ErrMarshal{
			Reasons: []string{
				"or element",
				"marshalling error: json: unsupported value: +Inf",
			},
		}.Error(),
	)
}

func TestSerializedCriteria_MarshallSuccesses(t *testing.T) {
	assert := testifyAssert.New(t)

	for _, testCase := range allMarshallUnmarshalTestCases {
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
	}
}
