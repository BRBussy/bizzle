package or

import (
	"errors"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
)

const Type = criterion.OrCriterionType

type Criterion struct {
	Criteria []criterion.Criterion
}

func (c Criterion) IsValid() error {
	if len(c.Criteria) == 0 {
		return errors.New("no criteria to or together")
	}
	return nil
}

func (c Criterion) Type() criterion.Type {
	return Type
}

func (c Criterion) ToFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	criteriaFilters := make([]map[string]interface{}, 0)
	for _, crit := range c.Criteria {
		criteriaFilters = append(criteriaFilters, crit.ToFilter())
	}
	filter["$or"] = criteriaFilters
	return filter
}
