package reasonInvalid

type Type string

func (t Type) String() string {
	return string(t)
}

const Unknown Type = "Unknown"
const Blank Type = "Blank"
const Nil Type = "Nil"
const DoesntExist Type = "DoesntExist"
const AlreadyExists Type = "AlreadyExists"
const Invalid Type = "Invalid"
const Duplicate Type = "Duplicate"
