package validator

import (
	"context"
)

type Validator interface {
	Validate(ctx context.Context, request ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	Token string `validate:"required"`
}

type ValidateResponse struct {
	MarshalledClaims []byte
}
