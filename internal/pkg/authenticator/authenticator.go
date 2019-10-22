package authenticator

type Authenticator interface {
	Login(*LoginRequest) (*LoginResponse, error)
}

const ServiceProvider = "Authenticator"
const LoginService = ServiceProvider + ".Login"

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type LoginResponse struct {
	JWT string
}
