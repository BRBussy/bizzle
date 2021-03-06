package operation

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
)

type And struct {
	Criteria []criterion.Criterion
}

func (a And) IsValid() error {
	if len(a.Criteria) == 0 {
		return criterion.ErrInvalid{Reasons: []string{"and operation criterion has an empty criterion array"}}
	}
	return nil
}

func (a And) Type() criterion.Type {
	return criterion.OperationAndCriterionType
}

func (a And) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	criteriaFilters := make([]map[string]interface{}, 0)
	for _, crit := range a.Criteria {
		criteriaFilters = append(criteriaFilters, crit.ToFilter())
	}
	filter["$and"] = criteriaFilters
	return filter
}

func (a And) ToJSON() (string, json.RawMessage, error) {
	return "", nil, criterion.ErrUnexpected{Reasons: []string{"and criterion to be marshalled during serialization"}}
}
