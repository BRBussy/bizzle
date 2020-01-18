package claims

type Type string

func (t Type) String() string {
	return string(t)
}

const LoginClaimsType Type = "Login"

type Claims interface {
	Type() Type              // Returns the Type of the claims
	ToJSON() ([]byte, error) // Returns json marshalled version of claims
	Expired() bool           // Returns whether or not the claims are expired
	ExpiryTime() int64       // Returns expiry time
}
