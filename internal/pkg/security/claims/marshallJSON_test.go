package claims

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type marshalJSONTestCase struct {
	name                     string
	marshalUnmarshalTestPair marshalUnmarshalTestPair
}

func TestSerialized_MarshalJSON(t *testing.T) {
	tests := []marshalJSONTestCase{
		{
			name:                     "loginClaims",
			marshalUnmarshalTestPair: loginClaimsTestPair,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				testSerialized := Serialized{
					Claims: tt.marshalUnmarshalTestPair.Claims,
				}
				jsonClaims, err := testSerialized.MarshalJSON()
				assert.Nil(t, err)
				assert.JSONEq(
					t,
					string(tt.marshalUnmarshalTestPair.Marshalled),
					string(jsonClaims),
				)
			},
		)
	}
}
