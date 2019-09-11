package criterion

type Type string

func (t Type) String() string {
	return string(t)
}

// criterion operations
const OrCriterionType Type = "OrCriterionType"

// basic criteria
const SubstringCriterionType Type = "SubstringCriterionType"
