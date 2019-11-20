package validator

import (
	"github.com/BRBussy/bizzle/internal/pkg/exercise/session"
	"github.com/BRBussy/bizzle/pkg/validate/reasonInvalid"
)

type Validator interface {
	ValidateForCreate(*ValidateForCreateRequest) (*ValidateForCreateResponse, error)
}

const ServiceProvider = "Session-Validator"
const ValidateForCreateService = ServiceProvider + ".ValidateForCreate"

type ValidateForCreateRequest struct {
	Session session.Session
}

type ValidateForCreateResponse struct {
	ReasonsInvalid reasonInvalid.ReasonsInvalid
}
