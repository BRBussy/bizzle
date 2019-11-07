package text

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
)

type Exact struct {
	Field string `json:"field"`
	Text  string `json:"text"`
}

func (e Exact) IsValid() error {

	reasonsInvalid := make([]string, 0)

	if e.Text == "" {
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
	return criterion.TextExactCriterionType
}

func (e Exact) ToFilter() map[string]interface{} {
	return map[string]interface{}{e.Field: e.Text}
}

func (e Exact) ToJSON() (string, json.RawMessage, error) {
	data, err := json.Marshal(struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}{
		Type: e.Type().String(),
		Text: e.Text,
	})
	return e.Field, data, err
}
