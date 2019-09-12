package criterion

import "github.com/stretchr/testify/suite"

type criterionTest struct {
	suite.Suite
}

func (t criterionTest) Test() {
	var testType Type = "testType"

	t.Equal(
		"testType",
		testType.String(),
	)

	t.Equal(
		"criterion is invalid: r1, r2",
		ErrInvalid{
			Reasons: []string{"r1", "r2"},
		}.Error(),
	)
}
