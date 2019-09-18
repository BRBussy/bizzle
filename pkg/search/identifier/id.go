package identifier

import (
	"encoding/json"
	"fmt"
)

type ID string

func (I ID) String() string {
	return string(I)
}

func (I ID) IsValid() error {
	if I == "" {
		return ErrInvalidIdentifier{Reasons: []string{"ID identifier is blank"}}
	}
	return nil
}

func (I ID) Type() Type {
	return IDIdentifierType
}

func (I ID) ToFilter() map[string]interface{} {
	return map[string]interface{}{"id": I.String()}
}

func (I ID) ToJSON() (map[string]json.RawMessage, error) {
	return map[string]json.RawMessage{
		"type": json.RawMessage(fmt.Sprintf(
			"\"%s\"",
			I.Type(),
		)),
		"id": json.RawMessage(fmt.Sprintf(
			"\"%s\"",
			I,
		)),
	}, nil
}
