package claims

func (s Serialized) MarshalJSON() ([]byte, error) {
	return s.Claims.ToJSON()
}
