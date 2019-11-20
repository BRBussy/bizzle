package exercise

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"go.mongodb.org/mongo-driver/bson"
)

type ArmCurl struct {
	ID        identifier.ID `json:"id" bson:"id"`
	SomeField string        `json:"someField" bson:"someField"`
}

func (c *ArmCurl) SetID(id identifier.ID) {
	c.ID = id
}

func (c *ArmCurl) Type() Type {
	return ArmCurlExerciseType
}

func (c *ArmCurl) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID        identifier.ID `json:"id"`
		Type      Type          `json:"type"`
		SomeField string        `json:"someField"`
	}{
		ID:        c.ID,
		Type:      c.Type(),
		SomeField: c.SomeField,
	})
}

func (c *ArmCurl) ToBSON() ([]byte, error) {
	return bson.Marshal(struct {
		ID        identifier.ID `bson:"id"`
		Type      Type          `bson:"type"`
		SomeField string        `bson:"someField"`
	}{
		ID:        c.ID,
		Type:      c.Type(),
		SomeField: c.SomeField,
	})
}
