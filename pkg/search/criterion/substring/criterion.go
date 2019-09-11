package substring

import (
	"errors"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	"strings"
)

const Type = criterion.Substring

type Criterion struct {
	Field string `json:"field"`
	Text  string `json:"text"`
}

func (c Criterion) IsValid() error {

	reasonsInvalid := make([]string, 0)

	if c.Text == "" {
		reasonsInvalid = append(reasonsInvalid, "text is blank")
	}

	if c.Field == "" {
		reasonsInvalid = append(reasonsInvalid, "field is blank")
	}

	if len(reasonsInvalid) > 0 {
		return errors.New(strings.Join(reasonsInvalid, "; "))
	}

	return nil
}

func (c Criterion) Type() criterion.Type {
	return Type
}

func (c Criterion) ToFilter() map[string]interface{} {
	return map[string]interface{}{
		"$regex":   ".*" + c.Text + ".*",
		"$options": "i",
	}
}
