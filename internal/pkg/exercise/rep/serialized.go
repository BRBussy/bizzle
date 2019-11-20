package rep

type Serialized struct {
	Rep Rep
}

func (s Serialized) MarshalBSON() ([]byte, error) {
	return s.Rep.ToBSON()
}

type bsonUnmarshalTypeHolder struct {
	Type Type `bson:"type"`
}
