package operation

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	suite.Run(t, new(andTest))
	suite.Run(t, new(orTest))
}
