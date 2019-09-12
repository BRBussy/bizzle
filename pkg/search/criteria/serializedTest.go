package criteria

import "github.com/stretchr/testify/suite"

type serializedTest struct {
	suite.Suite
}

func (t serializedTest) Test() {
	// testCase := "{\"name\":{\"type\":\"StringSubstring\",\"string\":\"sam\"}}"

	t.Equal(
		(&Serialized{}).UnmarshalJSON(nil),
		ErrInvalidSerializedCriteria{Reasons: []string{"json data nil"}},
	)
	t.Equal(
		(&Serialized{}).UnmarshalJSON([]byte("notValidJSON")),
		ErrUnmarshal{Reasons: []string{"json unmarshal", "invalid character 'o' in literal null (expecting 'u')"}},
	)
}
