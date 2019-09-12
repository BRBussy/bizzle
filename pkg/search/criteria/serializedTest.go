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
		(&Serialized{}).UnmarshalJSON(nil),
		ErrInvalidSerializedCriteria{Reasons: []string{"json criterion data is nil"}},
	)
	t.Equal(
		(&Serialized{}).UnmarshalJSON([]byte("notValidJSON")),
		ErrUnmarshal{Reasons: []string{"json unmarshal", "invalid character 'o' in literal null (expecting 'u')"}},
	)

	// test empty
	testEmpty := Serialized{}
	t.Equal(
		(&testEmpty).UnmarshalJSON([]byte("{}")),
		nil,
	)
	t.Equal(
		testEmpty.Criteria,
		make([]searchCriterion.Criterion, 0),
	)
}
