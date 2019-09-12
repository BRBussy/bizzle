package criteria

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	suite.Run(t, new(criteriaTest))
	suite.Run(t, new(serializedTest))
	suite.Run(t, new(errorTest))
}
