package wrapped

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/internal/pkg/errors"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	signInClaims "github.com/BRBussy/bizzle/internal/pkg/security/claims/signIn"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims/wrapped/exception"
	"net/http"
)

type Wrapped struct {
	Type  claims.Type     `json:"type"`
	Value json.RawMessage `json:"value"`
}

func Wrap(claimsToWrap claims.Claims) (Wrapped, error) {
	if claimsToWrap == nil {
		return Wrapped{}, exception.Invalid{Reasons: []string{"nil claimsToWrap provided"}}
	}

	marshalledValue, err := json.Marshal(claimsToWrap)
	if err != nil {
		return Wrapped{}, exception.Wrapping{Reasons: []string{"marshalling", err.Error()}}
	}
	return Wrapped{
		Type:  claimsToWrap.Type(),
		Value: marshalledValue,
	}, nil
}

func (wc Wrapped) Unwrap() (claims.Claims, error) {
	var result claims.Claims = nil

	switch wc.Type {
	case claims.SignIn:
		var unmarshalledClaims signInClaims.SignIn
		if err := json.Unmarshal(wc.Value, &unmarshalledClaims); err != nil {
			return nil, exception.Unwrapping{Reasons: []string{"unmarshalling", err.Error()}}
		}
		result = unmarshalledClaims

	default:
		return nil, exception.Invalid{Reasons: []string{"invalid type"}}
	}

	if result == nil {
		return nil, errors.ErrUnexpected{Reasons: []string{"identifier still nil"}}
	}

	// check for expiry
	if result.Expired() {
		return nil, exception.Invalid{Reasons: []string{"expired"}}
	}

	return result, nil
}

func UnwrapClaimsFromContext(r *http.Request) (claims.Claims, error) {
	wrapped, ok := r.Context().Value("wrappedClaims").(Wrapped)
	if !ok {
		return nil, exception.CouldNotParseFromContext{}
	}
	return wrapped.Unwrap()
}
