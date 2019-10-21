package adaptor

import (
	"net/http"
)

type ValidatorJSONRPCAdaptor struct {
	validator Validator
}

func (a *ValidatorJSONRPCAdaptor) Setup(validator Validator) *ValidatorJSONRPCAdaptor {
	return &ValidatorJSONRPCAdaptor{
		validator: validator,
	}
}

func (a *ValidatorJSONRPCAdaptor) ServiceName() string {
	return "TokenValidator"
}

type ValidateJSONRPCRequest struct {
	Token string `json:"token"`
}

type ValidateJSONRPCResponse struct {
	MarshalledClaims []byte `json:"marshalledClaims"`
}

func (a *ValidatorJSONRPCAdaptor) Validate(r *http.Request, request *ValidateJSONRPCRequest, response *ValidateJSONRPCResponse) error {
	validateResponse, err := a.validator.Validate(
		r.Context(),
		ValidateRequest{Token: request.Token},
	)
	if err != nil {
		return err
	}
	response.MarshalledClaims = validateResponse.MarshalledClaims
	return nil
}
