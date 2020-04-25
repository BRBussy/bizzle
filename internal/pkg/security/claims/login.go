package claims

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"time"
)

type Login struct {
	UserID         identifier.ID `validate:"required" json:"userID"`
	ExpirationTime int64         `validate:"required" json:"expirationTime"`
}

func (l Login) Type() Type {
	return LoginClaimsType
}

func (l Login) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type           Type          `json:"type"`
		UserID         identifier.ID `json:"userID"`
		ExpirationTime int64         `json:"expirationTime"`
	}{
		Type:           l.Type(),
		UserID:         l.UserID,
		ExpirationTime: l.ExpirationTime,
	})
}

func (l Login) Expired() bool {
	return time.Now().UTC().After(time.Unix(l.ExpirationTime, 0).UTC())
}

func (l Login) ExpiryTime() int64 {
	return l.ExpirationTime
}

func (l Login) ScopingID() identifier.ID {
	return l.UserID
}
