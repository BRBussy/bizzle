package criterion

type Type string

func (t Type) String() string {
	return string(t)
}

// criterion operators
const OROperator string = "$or"

// criterion operation types
const OperationOrCriterionType Type = "OperationOr"
const OperationAndCriterionType Type = "OperationAnd"

// string criterion types
const StringSubstringCriterionType Type = "StringSubstring"
const StringExactCriterionType Type = "StringExact"

// number criterion types
const NumberRangeCriterionType Type = "NumberRange"
const NumberExactCriterionType Type = "NumberExact"
