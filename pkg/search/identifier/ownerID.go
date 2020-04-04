package identifier

import (
	"encoding/json"
)

type OwnerID string

func (I OwnerID) String() string {
	return string(I)
}

func (I OwnerID) IsValid() error {
	if I == "" {
		return ErrInvalidIdentifier{Reasons: []string{"OwnerID identifier is blank"}}
	}
	return nil
}

func (I OwnerID) Type() Type {
	return OwnerIDIdentifierType
}

func (I OwnerID) ToFilter() map[string]interface{} {
	return map[string]interface{}{"ownerID": I.String()}
}

func (I OwnerID) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    Type   `json:"type"`
		OwnerID string `json:"ownerID"`
	}{
		Type:    I.Type(),
		OwnerID: I.String(),
	})
}
