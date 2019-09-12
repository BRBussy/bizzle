package operation

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
)

type Or struct {
	Criteria []criterion.Criterion
}

func (c Or) IsValid() error {
	if len(c.Criteria) == 0 {
		return criterion.ErrInvalid{Reasons: []string{"or operation criterion has an empty criterion array"}}
	}
	return nil
}

func (c Or) Type() criterion.Type {
	return criterion.OperationOrCriterionType
}

func (c Or) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	criteriaFilters := make([]map[string]interface{}, 0)
	for _, crit := range c.Criteria {
		criteriaFilters = append(criteriaFilters, crit.ToFilter())
	}
	filter["$or"] = criteriaFilters
	return filter
}
