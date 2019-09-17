package identifier

type ID struct {
	ID string `json:"id" bson:"id"`
}

// Type returns the Type of this Identifier
func (i ID) Type() Type { return IDIdentifierType }

// IsValid Determines and returns the validity of this Identifier
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

// ToFilter returns a query document for this identifier
func (i ID) ToFilter() map[string]interface{} {
	return map[string]interface{}{"id": i.ID}
}
