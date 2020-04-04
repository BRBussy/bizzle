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
		// try and get id field
		marshalledID, found := s.Serialized["id"]
		if !found {
			err := ErrInvalidSerializedIdentifier{Reasons: []string{"no id field"}}
			log.Error().Err(err)
			return err
		}
		// unmarshal ID identifier
		var idIdentifier ID
		if err := json.Unmarshal(marshalledID, &idIdentifier); err != nil {
			err = ErrUnmarshal{Reasons: []string{err.Error()}}
			log.Error().Err(err)
			return err
		}
		unmarshalledIdentifier = idIdentifier

	case OwnerIDIdentifierType:
		// try and get ownerID field
		marshalledID, found := s.Serialized["ownerID"]
		if !found {
			err := ErrInvalidSerializedIdentifier{Reasons: []string{"no ownerID field"}}
			log.Error().Err(err)
			return err
		}
		// unmarshal OwnerID identifier
		var ownerIDIdentifier OwnerID
		if err := json.Unmarshal(marshalledID, &ownerIDIdentifier); err != nil {
			err = ErrUnmarshal{Reasons: []string{err.Error()}}
			log.Error().Err(err)
			return err
		}
		unmarshalledIdentifier = ownerIDIdentifier

	case NameIdentifierType:
		// try and get name field
		marshalledID, found := s.Serialized["name"]
		if !found {
			err := ErrInvalidSerializedIdentifier{Reasons: []string{"no name field"}}
			log.Error().Err(err)
			return err
		}
		// unmarshal Name identifier
		var nameIdentifier Name
		if err := json.Unmarshal(marshalledID, &nameIdentifier); err != nil {
			err = ErrUnmarshal{Reasons: []string{err.Error()}}
			log.Error().Err(err)
			return err
		}
		unmarshalledIdentifier = nameIdentifier

	case EmailIdentifierType:
		// try and get email field
		marshalledID, found := s.Serialized["email"]
		if !found {
			err := ErrInvalidSerializedIdentifier{Reasons: []string{"no email field"}}
			log.Error().Err(err)
			return err
		}
		// unmarshal Email identifier
		var emailIdentifier Email
		if err := json.Unmarshal(marshalledID, &emailIdentifier); err != nil {
			err = ErrUnmarshal{Reasons: []string{err.Error()}}
			log.Error().Err(err)
			return err
		}
		unmarshalledIdentifier = emailIdentifier

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

	// set identifier
	s.Identifier = unmarshalledIdentifier
	return nil
}
