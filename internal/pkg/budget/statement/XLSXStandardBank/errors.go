package XLSXStandardBank

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
