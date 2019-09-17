package identifier

type ID struct {
	ID string `json:"id"`
}

// Returns IdentifierType of this Identifier
func (i ID) Type() Type { return IDIdentifierType }

// Determines and returns the validity of this Identifier
func (i ID) IsValid() error {
	reasonsInvalid := make([]string, 0)
	if i.ID == "" {
		reasonsInvalid = append(reasonsInvalid, "id is blank")
	}
	if len(reasonsInvalid) > 0 {
		return ErrInvalid{Reasons: reasonsInvalid}
	}
	return nil
}

func (i ID) ToFilter() map[string]interface{} {
	return map[string]interface{}{"id": i.ID}
}
