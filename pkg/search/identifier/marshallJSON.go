package identifier

import "encoding/json"

func (s *Serialized) MarshalJSON() ([]byte, error) {
	marshalledID, err := s.Identifier.ToJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(marshalledID)
}
