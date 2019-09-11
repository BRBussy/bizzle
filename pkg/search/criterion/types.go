package criterion

type Type string

func (t Type) String() string {
	return string(t)
}

// criterion operation types
const OperationOrCriterionType Type = "OperationOr"

// string criterion types
const StringSubstringCriterionType Type = "StringSubstring"
const StringExactCriterionType Type = "StringExact"
