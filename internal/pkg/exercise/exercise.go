package exercise

import "github.com/BRBussy/bizzle/pkg/search/identifier"

type Exercise interface {
	SetID(identifier.ID)     // Set ID on exercise
	Type() Type              // Returns the Type of the exercise
	ToJSON() ([]byte, error) // Returns json marshalled version of exercise
	ToBSON() ([]byte, error) // Returns bson marshalled version of exercise
}

type Type string

func (t Type) String() string {
	return string(t)
}

// arm exercises
const ArmCurlExerciseType Type = "ArmCurl"
