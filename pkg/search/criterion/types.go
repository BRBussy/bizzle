package criterion

type Type string

func (t Type) String() string {
	return string(t)
}

// criterion operations
const OrCriterionType Type = "OrCriterionType"

// basic criteria
const Substring Type = "Substring"

// list criteria
const ListText Type = "ListText"

// exact criteria
const ExactText Type = "ExactText"

// range criteria
const DateRange Type = "DateRange"
