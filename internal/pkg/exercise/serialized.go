package exercise

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

type Serialized struct {
	Exercise Exercise
}

func (s Serialized) MarshalBSON() ([]byte, error) {
	return s.Exercise.ToBSON()
}

type bsonUnmarshalTypeHolder struct {
	Type Type `bson:"type"`
}

func (s *Serialized) UnmarshalBSON(data []byte) error {
	// confirm that given data is not nil
	if data == nil {
		err := ErrInvalidSerializedExercise{Reasons: []string{"json exercise data is nil"}}
		log.Error().Err(err)
		return err
	}

	// unmarshal into type holder
	var th bsonUnmarshalTypeHolder
	if err := bson.Unmarshal(data, &th); err != nil {
		err = ErrUnmarshal{Reasons: []string{"json unmarshal into type holder", err.Error()}}
		log.Error().Err(err)
		return err
	}

	// unmarshal based on claims type
	var unmarshalledExercise Exercise
	switch th.Type {
	case ArmCurlExerciseType:
		var typedExercise ArmCurl
		if err := bson.Unmarshal(data, &typedExercise); err != nil {
			err = ErrUnmarshal{Reasons: []string{err.Error()}}
			log.Error().Err(err)
			return err
		}
		unmarshalledExercise = &typedExercise

	default:
		err := ErrInvalidSerializedExercise{
			Reasons: []string{
				"invalid type",
				th.Type.String(),
			},
		}
		log.Error().Err(err)
		return err
	}

	// set exercise
	s.Exercise = unmarshalledExercise
	return nil
}

func (s Serialized) MarshalJSON() ([]byte, error) {
	return s.Exercise.ToJSON()
}

type jsonUnmarshalTypeHolder struct {
	Type Type `json:"type"`
}

func (s *Serialized) UnmarshalJSON(data []byte) error {
	// confirm that given data is not nil
	if data == nil {
		err := ErrInvalidSerializedExercise{Reasons: []string{"json exercise data is nil"}}
		log.Error().Err(err)
		return err
	}

	// unmarshal into type holder
	var th jsonUnmarshalTypeHolder
	if err := json.Unmarshal(data, &th); err != nil {
		err = ErrUnmarshal{Reasons: []string{"json unmarshal into type holder", err.Error()}}
		log.Error().Err(err)
		return err
	}

	// unmarshal based on claims type
	var unmarshalledExercise Exercise
	switch th.Type {
	case ArmCurlExerciseType:
		var typedExercise ArmCurl
		if err := json.Unmarshal(data, &typedExercise); err != nil {
			err = ErrUnmarshal{Reasons: []string{err.Error()}}
			log.Error().Err(err)
			return err
		}
		unmarshalledExercise = &typedExercise

	default:
		err := ErrInvalidSerializedExercise{
			Reasons: []string{
				"invalid type",
				th.Type.String(),
			},
		}
		log.Error().Err(err)
		return err
	}

	// set exercise
	s.Exercise = unmarshalledExercise
	return nil
}
