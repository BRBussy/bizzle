package session

import (
	"github.com/BRBussy/bizzle/internal/pkg/exercise/rep"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Session struct {
	ID            string           `json:"id" bson:"id"`
	ExerciseReps  []rep.Serialized `json:"exerciseReps" bson:"exerciseReps"`
	DateTimeStart int64            `json:"dateTimeStart" bson:"dateTimeStart"`
	DateTimeEnd   int64            `json:"dateTimeEnd" bson:"dateTimeEnd"`
}

type ExerciseReps struct {
	ExerciseID identifier.ID `json:"exerciseID" bson:"exerciseID"`
}