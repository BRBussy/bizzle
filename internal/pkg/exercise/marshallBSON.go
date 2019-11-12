package exercise

func (s Serialized) MarshalBSON() ([]byte, error) {
	return s.Exercise.ToBSON()
}
