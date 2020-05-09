package admin

import (
	budgetEntryIgnored "github.com/BRBussy/bizzle/internal/pkg/budget/entry/ignored"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
)

type Admin interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
}

type CreateOneRequest struct {
	Claims  claims.Claims              `validate:"required"`
	Ignored budgetEntryIgnored.Ignored `validate:"required"`
}

type CreateOneResponse struct {
}
