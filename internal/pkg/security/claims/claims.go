package claims

import (
	"time"
)

type Type string

const SignIn Type = "SignIn"

type Claims interface {
	Type() Type
	Expired() bool
	TimeToExpiry() time.Duration
}
