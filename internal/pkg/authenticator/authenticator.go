package authenticator

import "github.com/BRBussy/bizzle/internal/pkg/security/claims"

type Authenticator interface {
	Login(*LoginRequest) (*LoginResponse, error)
	AuthenticateService(request *AuthenticateServiceRequest) (*AuthenticateServiceResponse, error)
}

const ServiceProvider = "Authenticator"
const LoginService = ServiceProvider + ".Login"
const AuthenticateServiceService = ServiceProvider + ".AuthenticateService"

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type LoginResponse struct {
	JWT string
}

type AuthenticateServiceRequest struct {
	Claims  claims.Claims `validate:"required"`
	Service string        `validate:"required"`
}

type AuthenticateServiceResponse struct {
}
