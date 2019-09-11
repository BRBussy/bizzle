package string

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
)

const Type = criterion.StringSubstringCriterionType

type Substring struct {
	Field string `json:"field"`
	Text  string `json:"text"`
}

func (c Substring) IsValid() error {

	reasonsInvalid := make([]string, 0)

	if c.Text == "" {
		reasonsInvalid = append(reasonsInvalid, "text is blank")
	}

	if c.Field == "" {
		reasonsInvalid = append(reasonsInvalid, "field is blank")
	}

	if len(reasonsInvalid) > 0 {
		return criterion.ErrInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (c Substring) Type() criterion.Type {
	return Type
}

func (c Substring) ToFilter() map[string]interface{} {
	return map[string]interface{}{
		"$regex":   ".*" + c.Text + ".*",
		"$options": "i",
	}
}
