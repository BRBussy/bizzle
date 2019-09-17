package criteria

import (
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestCriteriaToFilter(t *testing.T) {

}

func TestMarshalJSON(t *testing.T) {
	assert := testifyAssert.New(t)

	serializedCriterion, err := stringExact1.criteria.MarshalJSON()

}
