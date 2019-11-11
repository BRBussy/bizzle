package exercise

func (s Serialized) MarshalJSON() ([]byte, error) {
	return s.Exercise.ToJSON()
}
