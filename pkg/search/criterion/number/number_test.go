package number

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	suite.Run(t, new(exactTest))
	suite.Run(t, new(rangeTest))
}
