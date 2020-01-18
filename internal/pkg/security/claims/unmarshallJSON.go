package claims

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

type typeHolder struct {
	Type Type `json:"type"`
}

func (s *Serialized) UnmarshalJSON(data []byte) error {
	// confirm that given data is not nil
	if data == nil {
		err := ErrInvalidSerializedClaims{Reasons: []string{"json claims data is nil"}}
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
	var unmarshalledClaims Claims
	switch th.Type {
	case LoginClaimsType:
		var typedClaims Login
		if err := json.Unmarshal(data, &typedClaims); err != nil {
			err = ErrUnmarshal{Reasons: []string{err.Error()}}
			log.Error().Err(err)
			return err
		}
		unmarshalledClaims = typedClaims

	default:
		err := ErrInvalidSerializedClaims{
			Reasons: []string{
				"invalid type",
				th.Type.String(),
			},
		}
		log.Error().Err(err)
		return err
	}

	// set claims
	s.Claims = unmarshalledClaims
	return nil
}
