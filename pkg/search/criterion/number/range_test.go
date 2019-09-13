package number

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRange(t *testing.T) {
	assert := assert.New(t)

	// create a blank criterion
	testCriterion := Range{}

	// confirm is valid fails with field is blank
	assert.Equal(
		testCriterion.IsValid(),
		criterion.ErrInvalid{Reasons: []string{"field is blank"}},
	)

	// confirm that type returns correct type
	assert.Equal(
		criterion.NumberRangeCriterionType,
		testCriterion.Type(),
	)

	// populate field
	testCriterion.Field = "someField"

	// confirm is valid does not fail
	assert.Equal(
		nil,
		testCriterion.IsValid(),
	)

	// set start and end numbers
	testCriterion.Start.Number = 123.321
	testCriterion.End.Number = 321.123

	// test possible cases

	testCriterion.Start.Ignore = false
	testCriterion.Start.Inclusive = false
	testCriterion.End.Ignore = false
	testCriterion.End.Inclusive = false
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$gt": 123.321,
				"$lt": 321.123,
			},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = true
	testCriterion.Start.Inclusive = false
	testCriterion.End.Ignore = false
	testCriterion.End.Inclusive = false
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$lt": 321.123,
			},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = false
	testCriterion.Start.Inclusive = true
	testCriterion.End.Ignore = false
	testCriterion.End.Inclusive = false
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$gte": 123.321,
				"$lt":  321.123,
			},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = true
	testCriterion.Start.Inclusive = true
	testCriterion.End.Ignore = false
	testCriterion.End.Inclusive = false
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$lt": 321.123,
			},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = false
	testCriterion.Start.Inclusive = false
	testCriterion.End.Ignore = true
	testCriterion.End.Inclusive = false
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$gt": 123.321,
			},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = true
	testCriterion.Start.Inclusive = false
	testCriterion.End.Ignore = true
	testCriterion.End.Inclusive = false
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = false
	testCriterion.Start.Inclusive = true
	testCriterion.End.Ignore = true
	testCriterion.End.Inclusive = false
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$gte": 123.321,
			},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = true
	testCriterion.Start.Inclusive = true
	testCriterion.End.Ignore = true
	testCriterion.End.Inclusive = false
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = false
	testCriterion.Start.Inclusive = false
	testCriterion.End.Ignore = false
	testCriterion.End.Inclusive = true
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$gt":  123.321,
				"$lte": 321.123,
			},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = true
	testCriterion.Start.Inclusive = false
	testCriterion.End.Ignore = false
	testCriterion.End.Inclusive = true
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$lte": 321.123,
			},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = false
	testCriterion.Start.Inclusive = true
	testCriterion.End.Ignore = false
	testCriterion.End.Inclusive = true
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$gte": 123.321,
				"$lte": 321.123,
			},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = true
	testCriterion.Start.Inclusive = true
	testCriterion.End.Ignore = false
	testCriterion.End.Inclusive = true
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$lte": 321.123,
			},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = false
	testCriterion.Start.Inclusive = false
	testCriterion.End.Ignore = true
	testCriterion.End.Inclusive = true
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$gt": 123.321,
			},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = true
	testCriterion.Start.Inclusive = false
	testCriterion.End.Ignore = true
	testCriterion.End.Inclusive = true
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = false
	testCriterion.Start.Inclusive = true
	testCriterion.End.Ignore = true
	testCriterion.End.Inclusive = true
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{
				"$gte": 123.321,
			},
		},
		testCriterion.ToFilter(),
	)

	testCriterion.Start.Ignore = true
	testCriterion.Start.Inclusive = true
	testCriterion.End.Ignore = true
	testCriterion.End.Inclusive = true
	assert.Equal(
		map[string]interface{}{
			"someField": map[string]interface{}{},
		},
		testCriterion.ToFilter(),
	)
}
