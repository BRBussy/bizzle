package claims

import (
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"time"
)

type Login struct {
	UserID     identifier.ID `json:"userID"`
	ExpiryTime int64         `json:"expiryTime"`
}

func (l Login) Type() Type {
	return LoginClaims
}

func (l Login) Expired() bool {
	return time.Now().UTC().After(time.Unix(l.ExpiryTime, 0).UTC())
}
