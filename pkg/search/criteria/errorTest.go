package criteria

import "github.com/stretchr/testify/suite"

type errorTest struct {
	suite.Suite
}

func (t errorTest) Test() {
	t.Equal(
		ErrInvalidSerializedCriteria{Reasons: []string{"r1", "r2"}}.Error(),
		"serialized criteria is invalid: r1, r2",
	)
	t.Equal(
		ErrUnmarshal{Reasons: []string{"r1", "r2"}}.Error(),
		"unmarshalling error: r1, r2",
	)
}
