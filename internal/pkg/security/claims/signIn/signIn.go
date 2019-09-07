package signIn

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"time"
)

type SignIn struct {
	UserId         string `json:"userId"`
	IssueTime      int64  `json:"issueTime"`
	ExpirationTime int64  `json:"expirationTime"`
}

func (l SignIn) Type() claims.Type {
	return claims.SignIn
}

func (l SignIn) Expired() bool {
	return time.Now().UTC().After(time.Unix(l.ExpirationTime, 0).UTC())
}

func (l SignIn) TimeToExpiry() time.Duration {
	return time.Unix(l.ExpirationTime, 0).UTC().Sub(time.Now().UTC())
}
