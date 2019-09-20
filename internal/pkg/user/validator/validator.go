package validator

import (
	"github.com/BRBussy/bizzle/internal/pkg/user"
	"github.com/BRBussy/bizzle/pkg/validate/reasonInvalid"
)

type Validator interface {
	ValidateForCreate(request *ValidateForCreateRequest) (*ValidateForCreateResponse, error)
}

const ServiceProvider = "User-Validator"
const ValidateForCreateService = ServiceProvider + ".ValidateForCreate"

type ValidateForCreateRequest struct {
	User user.User
}

type ValidateForCreateResponse struct {
	ReasonsInvalid []reasonInvalid.ReasonInvalid
}
