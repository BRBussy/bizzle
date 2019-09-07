package secrets

import (
	"fmt"
	"strings"
)

type ErrCannotFindSecretFile struct {
	Path    string
	Reasons []string
}

func (e ErrCannotFindSecretFile) Error() string {
	return fmt.Sprintf(
		"cannot find secrets file at `%s`. %s",
		e.Path,
		strings.Join(e.Reasons, ", "),
	)
}
