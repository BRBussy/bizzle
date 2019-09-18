package criteria

import (
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	numberCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/number"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/string"
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestCriteriaToFilter(t *testing.T) {
	assert := testifyAssert.New(t)

	// zero length test
	var testCriteria Criteria = make([]searchCriterion.Criterion, 0)
	assert.Equal(
		0,
		len(testCriteria.ToFilter()),
	)

	stringExactCriterion := stringCriterion.Exact{
		Field:  "testStringExactField",
		String: "testStringExactContents",
	}
	stringSubstringCriterion := stringCriterion.Substring{
		Field:  "testStringSubstringField",
		String: "testStringSubstringContents",
	}
	numberExactCriterion := numberCriterion.Exact{
		Field:  "testNumberExactField",
		Number: 123.34,
	}

	// single length test
	testCriteria = []searchCriterion.Criterion{
		stringExactCriterion,
	}
	assert.Equal(
		stringExactCriterion.ToFilter(),
		testCriteria.ToFilter(),
	)

	// multiple test
	testCriteria = []searchCriterion.Criterion{
		stringExactCriterion,
		numberExactCriterion,
		stringSubstringCriterion,
	}
	filter := testCriteria.ToFilter()
	assert.Contains(
		filter,
		"$and",
		"filter should contain $and key",
	)
	andArray, found := filter["$and"]
	assert.Equal(
		true,
		found,
		"$and key should be found in filter",
	)
	andArrayTyped, ok := andArray.([]map[string]interface{})
	assert.Equal(
		true,
		ok,
		"should be able to infer type of and array",
	)
	assert.Equal(
		len(testCriteria),
		len(andArrayTyped),
		"and array should have correct number of elements",
	)
	assert.Contains(
		andArrayTyped,
		stringExactCriterion.ToFilter(),
		"string exact criterion should be in and array",
	)
	assert.Contains(
		andArrayTyped,
		numberExactCriterion.ToFilter(),
		"number exact criterion should be in and array",
	)
	assert.Contains(
		andArrayTyped,
		stringSubstringCriterion.ToFilter(),
		"string substring criterion should be in and array",
	)
}
