package claims

import (
	"context"
	"encoding/json"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
)

func ParseClaimsFromContext(ctx context.Context) (Claims, error) {
	// look for claims in context
	marshalledClaimsInterface := ctx.Value("Claims")
	if marshalledClaimsInterface == nil {
		return nil, ErrClaimsNotInContext{}
	}

	// try an cast claims to string
	marshalledClaims, ok := marshalledClaimsInterface.([]byte)
	if !ok {
		return nil, bizzleException.ErrUnexpected{Reasons: []string{"could not cast marshalled claims to []byte"}}
	}

	var serializedClaims Serialized
	if err := json.Unmarshal(marshalledClaims, &serializedClaims); err != nil {
		return nil, ErrUnmarshal{Reasons: []string{err.Error()}}
	}

	return serializedClaims.Claims, nil
}
