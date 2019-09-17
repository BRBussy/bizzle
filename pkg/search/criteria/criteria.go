package criteria

import (
	"encoding/json"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
)

// Criteria is a slice of criterion that can be converted into a filter
type Criteria []searchCriterion.Criterion

// ToFilter returns a filter combining all of the criterion in criteria
func (c Criteria) ToFilter() map[string]interface{} {
	filters := make([]map[string]interface{}, 0)
	for _, crit := range c {
		filters = append(filters, crit.ToFilter())
	}
	return map[string]interface{}{"$and": filters}
}

func (s *Serialized) MarshalJSON() ([]byte, error) {
	serializedCriteria := make(map[string]json.RawMessage)
	for _, criterion := range s.Criteria {

	}
}

func serialize(criterion searchCriterion.Criterion) (string, json.RawMessage, error) {
	return "", []byte{}, nil
}
