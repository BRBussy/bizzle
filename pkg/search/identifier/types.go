package identifier

type Type string

func (t Type) String() string {
	return string(t)
}

const IDIdentifierType Type = "ID"
const OwnerIDIdentifierType Type = "OwnerID"
const NameIdentifierType Type = "Name"
const EmailIdentifierType Type = "Email"
const NameVariantIdentifierType Type = "NameVariant"
