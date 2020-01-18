package operation

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
)

type Or struct {
	Criteria []criterion.Criterion
}

func (o Or) IsValid() error {
	if len(o.Criteria) == 0 {
		return criterion.ErrInvalid{Reasons: []string{"or operation criterion has an empty criterion array"}}
	}
	return nil
}

func (o Or) Type() criterion.Type {
	return criterion.OperationOrCriterionType
}

func (o Or) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	criteriaFilters := make([]map[string]interface{}, 0)
	for _, crit := range o.Criteria {
		criteriaFilters = append(criteriaFilters, crit.ToFilter())
	}
	filter["$or"] = criteriaFilters
	return filter
}

func (o Or) ToJSON() (string, json.RawMessage, error) {
	return "", nil, criterion.ErrUnexpected{Reasons: []string{"or criterion to be marshalled during serialization"}}
}
