package config

import (
	"fmt"

	"github.com/BRBussy/bizzle/pkg/validate/reasonInvalid"
)

// ErrInvalidConfig is a invalid budget config error
type ErrInvalidConfig struct {
	ReasonsInvalid reasonInvalid.ReasonsInvalid
}

func (e ErrInvalidConfig) Error() string {
	return fmt.Sprintf(
		"budget config invalid: %s",
		e.ReasonsInvalid.String(),
	)
}
