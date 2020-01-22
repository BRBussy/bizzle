package XLSXStandardBank

import (
	"fmt"
	"strings"
)

type ErrUnableToParseFile struct {
}

func (e ErrUnableToParseFile) Error() string {
	return "unable to parse excel file"
}

type ErrTransactionsSheetNotFound struct {
}

func (e ErrTransactionsSheetNotFound) Error() string {
	return "could not find transactions sheet"
}

type ErrSheetInvalid struct {
	Reasons []string
}

func (e ErrSheetInvalid) Error() string {
	return fmt.Sprintf(
		"sheet invalid: %s",
		strings.Join(e.Reasons, ", "),
	)
}
