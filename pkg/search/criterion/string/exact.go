package string

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
)

type Exact struct {
	Field  string `json:"field"`
	String string `json:"string"`
}

func (e Exact) IsValid() error {

	reasonsInvalid := make([]string, 0)

	if e.String == "" {
		reasonsInvalid = append(reasonsInvalid, "string is blank")
	}

	if e.Field == "" {
		reasonsInvalid = append(reasonsInvalid, "field is blank")
	}

	if len(reasonsInvalid) > 0 {
		return criterion.ErrInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (e Exact) Type() criterion.Type {
	return criterion.StringExactCriterionType
}

func (e Exact) ToFilter() map[string]interface{} {
	return map[string]interface{}{e.Field: e.String}
}

func (e Exact) ToJSON() (string, json.RawMessage, error) {
	data, err := json.Marshal(struct {
		Type   string `json:"type"`
		String string `json:"string"`
	}{
		Type:   e.Type().String(),
		String: e.String,
	})
	return e.Field, data, err
}
