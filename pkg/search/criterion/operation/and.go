package operation

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
)

type And struct {
	Criteria []criterion.Criterion
}

func (c And) IsValid() error {
	if len(c.Criteria) == 0 {
		return criterion.ErrInvalid{Reasons: []string{"and operation criterion has an empty criterion array"}}
	}
	return nil
}

func (c And) Type() criterion.Type {
	return criterion.OperationAndCriterionType
}

func (c And) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	criteriaFilters := make([]map[string]interface{}, 0)
	for _, crit := range c.Criteria {
		criteriaFilters = append(criteriaFilters, crit.ToFilter())
	}
	filter["$and"] = criteriaFilters
	return filter
}
