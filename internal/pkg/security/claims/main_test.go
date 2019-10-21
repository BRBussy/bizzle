package claims

import (
	"fmt"
)

type marshalUnmarshalTestPair struct {
	Marshalled []byte
	Claims     Claims
}

var loginClaimsTestPair = marshalUnmarshalTestPair{
	Marshalled: []byte(fmt.Sprintf(
		"{\"type\":\"%s\", \"expirationTime\":1570532806, \"userID\":\"1234\"}",
		LoginClaimsType,
	)),
	Claims: Login{
		UserID:         "1234",
		ExpirationTime: 1570532806,
	},
}
