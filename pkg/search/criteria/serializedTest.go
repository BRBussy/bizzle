package criteria

import (
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	"github.com/stretchr/testify/suite"
)

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

func (t serializedTest) TestSerializedCriteriaOROperator() {
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

func (t serializedTest) TestFieldCriterion() {
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
}
