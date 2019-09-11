package criterion

type Type string

// basic criteria
const Substring Type = "Substring"

// compound criteria
const Or Type = "Or"

// list criteria
const ListText Type = "ListText"

// exact criteria
const ExactText Type = "ExactText"

// range criteria
const DateRange Type = "DateRange"
