package session

import (
	"github.com/BRBussy/bizzle/internal/pkg/exercise/rep"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Session struct {
	ID            identifier.ID `json:"id" bson:"id"`
	Exercises     []Exercises   `json:"exercises" bson:"exercises"`
	DateTimeStart int64         `json:"dateTimeStart" bson:"dateTimeStart"`
	DateTimeEnd   int64         `json:"dateTimeEnd" bson:"dateTimeEnd"`
}

type Exercises struct {
	ExerciseID identifier.ID    `json:"exerciseID" bson:"exerciseID"`
	Reps       []rep.Serialized `json:"reps" bson:"reps"`
}
