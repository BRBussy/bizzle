package text

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
)

type List struct {
	Field string   `json:"field"`
	List  []string `json:"list"`
}

func (l List) IsValid() error {
	reasonsInvalid := make([]string, 0)
	if l.Field == "" {
		reasonsInvalid = append(reasonsInvalid, "field is blank")
	}
	if len(l.List) == 0 {
		reasonsInvalid = append(reasonsInvalid, "list is empty")
	}
	if len(reasonsInvalid) > 0 {
		return criterion.ErrInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (l List) Type() criterion.Type {
	return criterion.TextListCriterionType
}

func (l List) ToFilter() map[string]interface{} {
	return map[string]interface{}{
		l.Field: map[string]interface{}{
			"$in": l.List,
		},
	}
}

func (l List) ToJSON() (string, json.RawMessage, error) {
	data, err := json.Marshal(struct {
		Type string   `json:"type"`
		List []string `json:"list"`
	}{
		Type: l.Type().String(),
		List: l.List,
	})
	return l.Field, data, err
}
