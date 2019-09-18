package criteria

import (
	"encoding/json"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
)

type Serialized struct {
	Serialized map[string]json.RawMessage
	Criteria   []searchCriterion.Criterion
}
