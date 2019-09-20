package identifier

type Type string

func (t Type) String() string {
	return string(t)
}

const IDIdentifierType Type = "ID"
const NameIdentifierType Type = "Name"
const EmailIdentifierType Type = "Email"
