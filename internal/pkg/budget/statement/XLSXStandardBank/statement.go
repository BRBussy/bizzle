package XLSXStandardBank

type ColumnHeader string

func (c ColumnHeader) String() string {
	return string(c)
}

const DateColumnHeader ColumnHeader = "Date"
const DescriptionColumnHeader ColumnHeader = "Description"
const InColumnHeader ColumnHeader = "In (R)"
const OutColumnHeader ColumnHeader = "Out (R)"
const BankFeesColumnHeader ColumnHeader = "Bank fees (R)"
const BalanceColumnHeader ColumnHeader = "Balance (R)"

var RequiredColumnHeaders = []ColumnHeader{
	DateColumnHeader,
	DescriptionColumnHeader,
	InColumnHeader,
	OutColumnHeader,
	BalanceColumnHeader,
}
