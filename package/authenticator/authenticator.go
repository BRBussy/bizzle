package authenticator

type Authenticator interface {
	SignUp(*SignUpRequest) (*SignUpResponse, error)
}

const ServiceProvider = "Authenticator"
const SignUpService = ServiceProvider + ".SignUp"

type SignUpRequest struct {
}

type SignUpResponse struct {
	Msg string
}
