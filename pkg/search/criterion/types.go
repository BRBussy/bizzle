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
const TextSubstringCriterionType Type = "TextSubstring"
const TextExactCriterionType Type = "TextExact"
const TextListCriterionType Type = "TextList"

// number criterion types
const NumberRangeCriterionType Type = "NumberRange"
const NumberExactCriterionType Type = "NumberExact"

// dateTime criterion types
const DateTimeRangeCriterionType Type = "DateTimeRange"
