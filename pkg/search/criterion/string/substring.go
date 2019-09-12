package string

import (
	"github.com/BRBussy/bizzle/pkg/search/criterion"
)

type Substring struct {
	Field  string `json:"field"`
	String string `json:"string"`
}

func (s Substring) IsValid() error {

	reasonsInvalid := make([]string, 0)

	if s.String == "" {
		reasonsInvalid = append(reasonsInvalid, "text is blank")
	}

	if s.Field == "" {
		reasonsInvalid = append(reasonsInvalid, "field is blank")
	}

	if len(reasonsInvalid) > 0 {
		return criterion.ErrInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (s Substring) Type() criterion.Type {
	return criterion.StringSubstringCriterionType
}

func (s Substring) ToFilter() map[string]interface{} {
	return map[string]interface{}{
		s.Field: map[string]interface{}{
			"$regex":   ".*" + s.String + ".*",
			"$options": "i",
		},
	}
}
