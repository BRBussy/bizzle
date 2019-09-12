package criteria

import (
	"fmt"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	operationCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/operation"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/string"
	"github.com/stretchr/testify/suite"
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

type serializedTest struct {
	suite.Suite
}

func (t serializedTest) TestSerializedCriteriaInvalidInput() {
	// test fringe invalid json inputs
	t.Equal(
		ErrInvalidSerializedCriteria{Reasons: []string{"json criterion data is nil"}},
		(&Serialized{}).UnmarshalJSON(nil),
	)
	t.Equal(
		ErrUnmarshal{Reasons: []string{"json unmarshal", "invalid character 'o' in literal null (expecting 'u')"}},
		(&Serialized{}).UnmarshalJSON([]byte("notValidJSON")),
	)

	// test empty
	testEmpty := Serialized{}
	t.Equal(
		nil,
		(&testEmpty).UnmarshalJSON([]byte("{}")),
	)
	t.Equal(
		make([]searchCriterion.Criterion, 0),
		testEmpty.Criteria,
	)
}

func (t serializedTest) TestSerializedCriteriaOROperatorFailures() {
	// invalid value provided for $or operator
	t.Equal(
		ErrUnmarshal{Reasons: []string{
			"or array unmarshal",
			"json: cannot unmarshal object into Go value of type criteria.jsonObjectArray",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"$or\":{}}")),
	)
	t.Equal(
		ErrUnmarshal{Reasons: []string{
			"or array unmarshal",
			"json: cannot unmarshal number into Go value of type criteria.jsonObjectArray",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"$or\":6}")),
	)
	t.Equal(
		ErrUnmarshal{Reasons: []string{
			"or array unmarshal",
			"json: cannot unmarshal string into Go value of type criteria.jsonObjectArray",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"$or\":\"string\"}")),
	)
	// empty array
	t.Equal(
		ErrInvalidSerializedCriteria{Reasons: []string{
			searchCriterion.ErrInvalid{Reasons: []string{"or operation criterion has an empty criterion array"}}.Error(),
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"$or\":[]}")),
	)
}

func (t serializedTest) TestSerializedCriteriaFieldCriterionFailures() {
	// invalid input for field criterion
	t.Equal(
		ErrUnmarshal{Reasons: []string{
			"criterion object unmarshal",
			"json: cannot unmarshal array into Go value of type criteria.typeHolder",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"someField\":[]}")),
	)
	t.Equal(
		ErrUnmarshal{Reasons: []string{
			"criterion object unmarshal",
			"json: cannot unmarshal number into Go value of type criteria.typeHolder",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"someField\":5}")),
	)
	t.Equal(
		ErrUnmarshal{Reasons: []string{
			"criterion object unmarshal",
			"json: cannot unmarshal string into Go value of type criteria.typeHolder",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"someField\":\"string\"}")),
	)
	// invalid type for field criterion
	t.Equal(
		ErrUnmarshal{Reasons: []string{
			"criterion object unmarshal",
			"json: cannot unmarshal number into Go struct field typeHolder.type of type criterion.Type",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"someField\":{\"type\":4}}")),
	)
	// invalid value for field criterion
	t.Equal(
		ErrInvalidSerializedCriteria{Reasons: []string{
			"invalid field criterion type",
			"invalidType",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"someField\":{\"type\":\"invalidType\"}}")),
	)
}

func (t serializedTest) TestSerializedCriteriaStringSubstringCriterion() {
	testSubstringCriterion := Serialized{}
	t.Equal(
		nil,
		(&testSubstringCriterion).UnmarshalJSON(stringSubstring1.serializedCriterion),
	)
	t.Equal(
		[]searchCriterion.Criterion{
			stringSubstring1.criterion,
		},
		testSubstringCriterion.Criteria,
	)
}

func (t serializedTest) TestSerializedCriteriaStringExactCriterion() {
	testSubstringCriterion := Serialized{}
	t.Equal(
		nil,
		(&testSubstringCriterion).UnmarshalJSON(stringExact1.serializedCriterion),
	)
	t.Equal(
		[]searchCriterion.Criterion{
			stringExact1.criterion,
		},
		testSubstringCriterion.Criteria,
	)
}

func (t serializedTest) TestSerializedCriteriaOrCriterion() {
	serializedValue := []byte(fmt.Sprintf(
		"{\"$or\":[%s,%s]}",
		stringSubstring1.serializedCriterion,
		stringExact1.serializedCriterion,
	))

	testOrCriterion := Serialized{}

	t.Equal(
		nil,
		(&testOrCriterion).UnmarshalJSON(serializedValue),
	)
	t.Equal(
		[]searchCriterion.Criterion{
			operationCriterion.Or{
				Criteria: []searchCriterion.Criterion{
					stringSubstring1.criterion,
					stringExact1.criterion,
				},
			},
		},
		testOrCriterion.Criteria,
	)
}
