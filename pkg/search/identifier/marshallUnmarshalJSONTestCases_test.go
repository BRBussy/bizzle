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

var allTestCases = []TestCase{
	idIdentifierTestCase0,
}
