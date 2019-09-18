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
		err = ErrUnmarshal{Reasons: []string{"json unmarshal", err.Error()}}
		log.Error().Err(err)
		return err
	}

	return nil
}
