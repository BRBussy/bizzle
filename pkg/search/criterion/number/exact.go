package number

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
)

type Exact struct {
	Field  string  `json:"field"`
	Number float64 `json:"number"`
}

func (e Exact) IsValid() error {

	reasonsInvalid := make([]string, 0)

	if e.Field == "" {
		reasonsInvalid = append(reasonsInvalid, "field is blank")
	}

	if len(reasonsInvalid) > 0 {
		return criterion.ErrInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (e Exact) Type() criterion.Type {
	return criterion.NumberExactCriterionType
}

func (e Exact) ToFilter() map[string]interface{} {
	return map[string]interface{}{e.Field: e.Number}
}

func (e Exact) ToJSON() (string, json.RawMessage, error) {
	data, err := json.Marshal(struct {
		Type   string  `json:"type"`
		Number float64 `json:"number"`
	}{
		Type:   e.Type().String(),
		Number: e.Number,
	})
	return e.Field, data, err
}
