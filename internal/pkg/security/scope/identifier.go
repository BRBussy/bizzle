package scope

import "github.com/BRBussy/bizzle/pkg/search/identifier"

type Identifier struct {

}

func (i Identifier) Type() {
	return identifier.Type("ScopeIdentifier")
}

func (i Identifier) IsValid() error {
	return nil
}

ToFilter() map[string]interface{}
ToJSON() ([]byte, error)

