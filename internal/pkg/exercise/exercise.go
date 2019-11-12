package exercise

import "github.com/BRBussy/bizzle/pkg/search/identifier"

type Exercise interface {
	SetID(identifier.ID)     // Set ID on exercise
	Type() Type              // Returns the Type of the exercise
	ToJSON() ([]byte, error) // Returns json marshalled version of exercise
}
