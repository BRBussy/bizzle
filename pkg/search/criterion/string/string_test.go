package string

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestExact(t *testing.T) {
	suite.Run(t, new(exactTest))
	suite.Run(t, new(substringTest))
}
