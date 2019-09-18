package identifier

import "encoding/json"

type Identifier interface {
	IsValid() error                              // Returns the validity of the Identifier
	Type() Type                                  // Returns the Type of the Identifier
	ToFilter() map[string]interface{}            // Returns a map filter to use to query the databases
	ToJSON() (map[string]json.RawMessage, error) // Returns js 'pojo' of identifier
}
