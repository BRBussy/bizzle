package arm

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/internal/pkg/exercise"
)

type Curl struct {
	SomeField string `json:"someField" bson:"someField"`
}

func (c Curl) Type() exercise.Type {
	return exercise.ArmCurlExerciseType
}

func (c Curl) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      exercise.Type `json:"type"`
		SomeField string        `json:"someField"`
	}{
		Type:      c.Type(),
		SomeField: c.SomeField,
	})
}
