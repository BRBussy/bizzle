package identifier

import (
	"encoding/json"
)

type NameVariant struct {
	Name    string
	Variant string
}

func (n NameVariant) IsValid() error {
	reasonsInvalid := make([]string, 0)
	if n.Name == "" {
		reasonsInvalid = append(reasonsInvalid, "name is blank")
	}
	if n.Variant == "" {
		reasonsInvalid = append(reasonsInvalid, "variant is blank")
	}

	if len(reasonsInvalid) > 0 {
		return ErrInvalidIdentifier{Reasons: reasonsInvalid}
	}
	return nil
}

func (n NameVariant) Type() Type {
	return NameVariantIdentifierType
}

func (n NameVariant) ToFilter() map[string]interface{} {
	return map[string]interface{}{
		"name":    n.Name,
		"variant": n.Variant,
	}
}

func (n NameVariant) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    Type   `json:"type"`
		Name    string `json:"name"`
		Variant string `json:"variant"`
	}{
		Type:    n.Type(),
		Name:    n.Name,
		Variant: n.Variant,
	})
}
