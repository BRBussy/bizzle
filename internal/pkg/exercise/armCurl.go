package exercise

import (
	"encoding/json"
)

type ArmCurl struct {
	SomeField string `json:"someField" bson:"someField"`
}

func (c ArmCurl) Type() Type {
	return ArmCurlExerciseType
}

func (c ArmCurl) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      Type   `json:"type"`
		SomeField string `json:"someField"`
	}{
		Type:      c.Type(),
		SomeField: c.SomeField,
	})
}
