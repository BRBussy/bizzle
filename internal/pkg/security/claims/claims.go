package claims

type Type string

const LoginClaims Type = "Login"

type Claims interface {
	Type() Type
	Expired() bool
}
