package rep

import "github.com/BRBussy/bizzle/pkg/search/identifier"

type Rep interface {
	SetID(identifier.ID)     // Set ID on rep
	Type() Type              // Returns the Type of the rep
	ToJSON() ([]byte, error) // Returns json marshalled version of rep
	ToBSON() ([]byte, error) // Returns bson marshalled version of rep
}

type Type string

func (t Type) String() string {
	return string(t)
}

const DurationRepType Type = "Duration"
