package scope

import (
	"encoding/json"

	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

// ScopedIdentifier is used to apply an owner ID filter to identifier
type ScopedIdentifier struct {
	IdentifierToScope identifier.Identifier
	OwnerID           identifier.ID
}

// Type returns the identifer.Type of the identifier
func (i ScopedIdentifier) Type() identifier.Type {
	return identifier.Type("ScopedIdentifier")
}

// IsValid determines if the identifier is valid
func (i ScopedIdentifier) IsValid() error {
	return nil
}

// ToFilter converts the identifier into a filter for mongo db queries
func (i ScopedIdentifier) ToFilter() map[string]interface{} {
	scopedFilter := map[string]interface{}{
		"ownerID": i.OwnerID,
	}
	for key, val := range i.IdentifierToScope.ToFilter() {
		scopedFilter[key] = val
	}
	return scopedFilter
}

// ToJSON returns the identifier JSON marshalled
func (i ScopedIdentifier) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type              identifier.Type       `json:"type`
		IdentifierToScope identifier.Serialized `json:"identifierToScope`
	}{
		Type: i.Type(),
		IdentifierToScope: identifier.Serialized{
			Identifier: i.IdentifierToScope,
		},
	})
}
