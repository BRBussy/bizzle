package identifier

import (
	"encoding/json"
)

type Name string

func (n Name) String() string {
	return string(n)
}

func (n Name) IsValid() error {
	if n == "" {
		return ErrInvalidIdentifier{Reasons: []string{"Name identifier is blank"}}
	}
	return nil
}

func (n Name) Type() Type {
	return NameIdentifierType
}

func (n Name) ToFilter() map[string]interface{} {
	return map[string]interface{}{"name": n.String()}
}

func (n Name) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type Type   `json:"type"`
		Name string `json:"name"`
	}{
		Type: n.Type(),
		Name: n.String(),
	})
}
