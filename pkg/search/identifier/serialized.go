package identifier

import (
	"encoding/json"
)

type Serialized struct {
	Serialized map[string]json.RawMessage
	Identifier Identifier
}
