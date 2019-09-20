package identifier

import (
	"encoding/json"
)

type Email string

func (n Email) String() string {
	return string(n)
}

func (n Email) IsValid() error {
	if n == "" {
		return ErrInvalidIdentifier{Reasons: []string{"Email identifier is blank"}}
	}
	return nil
}

func (n Email) Type() Type {
	return EmailIdentifierType
}

func (n Email) ToFilter() map[string]interface{} {
	return map[string]interface{}{"email": n.String()}
}

func (n Email) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  Type   `json:"type"`
		Email string `json:"email"`
	}{
		Type:  n.Type(),
		Email: n.String(),
	})
}
