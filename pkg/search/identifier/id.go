package identifier

type ID string

func (I ID) String() string {
	return string(I)
}

func (I ID) IsValid() error {
	if I == "" {
		return ErrInvalidIdentifier{Reasons: []string{"ID identifier is blank"}}
	}
	return nil
}

func (I ID) Type() Type {
	return IDIdentifierType
}

func (I ID) ToFilter() map[string]interface{} {
	return map[string]interface{}{"id": I.String()}
}
