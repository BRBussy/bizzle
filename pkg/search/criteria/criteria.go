package criteria

import (
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
)

// Criteria is a slice of criterion that can be converted into a filter
type Criteria []searchCriterion.Criterion

// ToFilter returns a filter combining all of the criterion in criteria
func (c Criteria) ToFilter() map[string]interface{} {
	filters := make([]map[string]interface{}, 0)

	if len(c) == 0 {
		// if the criteria array has no entries return blank object
		return make(map[string]interface{})
	} else if len(c) == 1 {
		// if there is only one entry, return just that single filter
		return c[0].ToFilter()
	}

	// otherwise return a combined and filter
	for _, crit := range c {
		filters = append(filters, crit.ToFilter())
	}
	return map[string]interface{}{"$and": filters}
}
