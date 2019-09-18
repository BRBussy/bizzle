package identifier

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

func (s *Serialized) UnmarshalJSON(data []byte) error {
	// confirm that given data is not nil
	if data == nil {
		err := ErrInvalidSerializedIdentifier{Reasons: []string{"json identifier data is nil"}}
		log.Error().Err(err)
		return err
	}

	// unmarshal into serialized section of Serialized
	if err := json.Unmarshal(data, &s.Serialized); err != nil {
		err = ErrUnmarshal{Reasons: []string{"json unmarshal id object", err.Error()}}
		log.Error().Err(err)
		return err
	}

	// try and get type field
	marshalledIdType, found := s.Serialized["type"]
	if !found {
		err := ErrInvalidSerializedIdentifier{Reasons: []string{"no type field"}}
		log.Error().Err(err)
		return err
	}

	// unmarshal id type
	var idType Type
	if err := json.Unmarshal(marshalledIdType, &idType); err != nil {
		err = ErrUnmarshal{Reasons: []string{"json unmarshal id type", err.Error()}}
		log.Error().Err(err)
		return err
	}

	// unmarshal based on identifier type
	var unmarshalledIdentifier Identifier
	switch idType {
	case IDIdentifierType:
		var idIdentifier ID
		if err := json.Unmarshal(data, &idIdentifier); err != nil {
			err = ErrUnmarshal{Reasons: []string{err.Error()}}
			log.Error().Err(err)
			return err
		}
		unmarshalledIdentifier = idIdentifier

	default:
		err := ErrInvalidSerializedIdentifier{
			Reasons: []string{
				"invalid type",
				string(idType),
			},
		}
		log.Error().Err(err)
		return err
	}

	// check validity
	if err := unmarshalledIdentifier.IsValid(); err != nil {
		return err
	}

	return nil
}
