package criteria

import (
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	"github.com/stretchr/testify/suite"
)

type serializedTest struct {
	suite.Suite
}

func (t serializedTest) Test() {
	// testCase := "{\"name\":{\"type\":\"StringSubstring\",\"string\":\"sam\"}}"

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

	t.Equal(
		ErrUnmarshal{Reasons: []string{
			"or array",
		}},
		(&Serialized{}).UnmarshalJSON([]byte("{\"$or\":[]}")),
	)
}
