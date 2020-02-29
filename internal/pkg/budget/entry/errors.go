package entry

import (
	"fmt"

	"github.com/BRBussy/bizzle/pkg/validate/reasonInvalid"
)

// ErrInvalidEntry is a invalid budget entry error
type ErrInvalidEntry struct {
	ReasonsInvalid reasonInvalid.ReasonsInvalid
}

func (e ErrInvalidEntry) Error() string {
	return fmt.Sprintf(
		"budget entry invalid: %s",
		e.ReasonsInvalid.String(),
	)
}
