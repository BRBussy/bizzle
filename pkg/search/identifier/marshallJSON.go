package identifier

func (s *Serialized) MarshalJSON() ([]byte, error) {
	return s.Identifier.ToJSON()
}
