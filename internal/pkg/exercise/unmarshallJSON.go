package exercise

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/internal/pkg/exercise/arm"
	"github.com/rs/zerolog/log"
)

type typeHolder struct {
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
	var th typeHolder
	if err := json.Unmarshal(data, &th); err != nil {
		err = ErrUnmarshal{Reasons: []string{"json unmarshal into type holder", err.Error()}}
		log.Error().Err(err)
		return err
	}

	// unmarshal based on claims type
	var unmarshalledExercise Exercise
	switch th.Type {
	case ArmCurlExerciseType:
		var typedExercise arm.Curl
		if err := json.Unmarshal(data, &typedExercise); err != nil {
			err = ErrUnmarshal{Reasons: []string{err.Error()}}
			log.Error().Err(err)
			return err
		}
		unmarshalledExercise = typedExercise

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
