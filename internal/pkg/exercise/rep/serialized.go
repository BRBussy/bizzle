package rep

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

type Serialized struct {
	Rep Rep
}

func (s Serialized) MarshalBSON() ([]byte, error) {
	return s.Rep.ToBSON()
}

type bsonUnmarshalTypeHolder struct {
	Type Type `bson:"type"`
}

func (s *Serialized) UnmarshalBSON(data []byte) error {
	// confirm that given data is not nil
	if data == nil {
		err := ErrInvalidSerializedRep{Reasons: []string{"json rep data is nil"}}
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
	var unmarshalledRep Rep
	switch th.Type {
	case DurationRepType:
		var typedRep Duration
		if err := bson.Unmarshal(data, &typedRep); err != nil {
			err = ErrUnmarshal{Reasons: []string{err.Error()}}
			log.Error().Err(err)
			return err
		}
		unmarshalledRep = &typedRep

	default:
		err := ErrInvalidSerializedRep{
			Reasons: []string{
				"invalid type",
				th.Type.String(),
			},
		}
		log.Error().Err(err)
		return err
	}

	s.Rep = unmarshalledRep
	return nil
}

func (s Serialized) MarshalJSON() ([]byte, error) {
	return s.Rep.ToJSON()
}

type jsonUnmarshalTypeHolder struct {
	Type Type `json:"type"`
}

func (s *Serialized) UnmarshalJSON(data []byte) error {
	// confirm that given data is not nil
	if data == nil {
		err := ErrInvalidSerializedRep{Reasons: []string{"json exercise data is nil"}}
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
	var unmarshalledRep Rep
	switch th.Type {
	case DurationRepType:
		var typedRep Duration
		if err := json.Unmarshal(data, &typedRep); err != nil {
			err = ErrUnmarshal{Reasons: []string{err.Error()}}
			log.Error().Err(err)
			return err
		}
		unmarshalledRep = &typedRep

	default:
		err := ErrInvalidSerializedRep{
			Reasons: []string{
				"invalid type",
				th.Type.String(),
			},
		}
		log.Error().Err(err)
		return err
	}

	// set exercise
	s.Rep = unmarshalledRep
	return nil
}
