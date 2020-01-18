package identifier

import "fmt"

type TestCase struct {
	id                   string
	SerializedIdentifier []byte
	Identifier           Identifier
}

var idIdentifierTestCase0 = TestCase{
	id: "idIdentifierTestCase0",
	SerializedIdentifier: []byte(fmt.Sprintf(
		"{\"type\":\"%s\",\"id\":\"1234\"}",
		IDIdentifierType,
	)),
	Identifier: ID("1234"),
}

var nameIdentifierTestCase0 = TestCase{
	id: "nameIdentifierTestCase0",
	SerializedIdentifier: []byte(fmt.Sprintf(
		"{\"type\":\"%s\",\"name\":\"bob\"}",
		NameIdentifierType,
	)),
	Identifier: Name("bob"),
}

var emailIdentifierTestCase0 = TestCase{
	id: "emailIdentifierTestCase0",
	SerializedIdentifier: []byte(fmt.Sprintf(
		"{\"type\":\"%s\",\"email\":\"bob@gmail.com\"}",
		EmailIdentifierType,
	)),
	Identifier: Email("bob@gmail.com"),
}

var allTestCases = []TestCase{
	idIdentifierTestCase0,
	nameIdentifierTestCase0,
	emailIdentifierTestCase0,
}
