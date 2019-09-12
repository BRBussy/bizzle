package criteria

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
)

// Criteria
type Criteria interface {
	IsValid() error                   // Returns the validity of the Criterion
	ToFilter() map[string]interface{} // Returns a map filter to use to query the databases
}

type SerializedCriteria struct {
	Serialized map[string]interface{}
}

func (s *SerializedCriteria) UnmarshalJSON(data []byte) error {
	// unmarshal into serialized section of SerializedCriteria
	if err := json.Unmarshal(data, &s.Serialized); err != nil {
		log.Error().Err(err).Msg("unmarshalling wrapped criterion")
		return errors.New("unmarshalling failed: " + err.Error())
	}

	return nil
}
