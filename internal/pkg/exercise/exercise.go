package exercise

import "github.com/BRBussy/bizzle/pkg/search/identifier"

type Exercise struct {
	ID          identifier.ID `json:"id" bson:"id"`
	Name        string        `json:"name" bson:"name"`
	MuscleGroup MuscleGroup   `json:"muscleGroup" bson:"muscleGroup"`
	Variant     string        `json:"variant" bson:"variant"`
	Description string        `json:"description" bson:"description"`
}

type MuscleGroup string

func (m MuscleGroup) String() string {
	return string(m)
}

const BicepsMuscleGroup MuscleGroup = "Biceps"
