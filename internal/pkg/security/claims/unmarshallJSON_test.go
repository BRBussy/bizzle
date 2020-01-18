package claims

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type unmarshalJSONTestCase struct {
	name                     string
	marshalUnmarshalTestPair marshalUnmarshalTestPair
	check                    func(t *testing.T, pair marshalUnmarshalTestPair, err error, serialized Serialized)
}

func TestSerialized_UnmarshalJSON(t *testing.T) {
	tests := []unmarshalJSONTestCase{
		{
			name: "nil data",
			marshalUnmarshalTestPair: marshalUnmarshalTestPair{
				Marshalled: nil,
			},
			check: func(t *testing.T, pair marshalUnmarshalTestPair, err error, testSerialized Serialized) {
				assert.Nil(t, testSerialized.Claims)
				assert.IsType(t, ErrInvalidSerializedClaims{}, err)
				assert.EqualError(t, err, "invalid serialized claims: json claims data is nil")
			},
		},
		{
			name: "unmarshal error - invalid json",
			marshalUnmarshalTestPair: marshalUnmarshalTestPair{
				Marshalled: []byte("notValidJson"),
			},
			check: func(t *testing.T, pair marshalUnmarshalTestPair, err error, testSerialized Serialized) {
				assert.Nil(t, testSerialized.Claims)
				assert.IsType(t, ErrUnmarshal{}, err)
				assert.EqualError(t, err, "unmarshalling error: json unmarshal into type holder, invalid character 'o' in literal null (expecting 'u')")
			},
		},
		{
			name: "unmarshal error - invalid type",
			marshalUnmarshalTestPair: marshalUnmarshalTestPair{
				Marshalled: []byte("{\"type\":1234, \"userID\":\"1234\"}"),
			},
			check: func(t *testing.T, pair marshalUnmarshalTestPair, err error, testSerialized Serialized) {
				assert.Nil(t, testSerialized.Claims)
				assert.IsType(t, ErrUnmarshal{}, err)
				assert.EqualError(t, err, "unmarshalling error: json unmarshal into type holder, json: cannot unmarshal number into Go struct field typeHolder.type of type claims.Type")
			},
		},
		{
			name: "invalid claims - invalid type",
			marshalUnmarshalTestPair: marshalUnmarshalTestPair{
				Marshalled: []byte("{\"type\":\"notCorrect\", \"userID\":\"1234\"}"),
			},
			check: func(t *testing.T, pair marshalUnmarshalTestPair, err error, testSerialized Serialized) {
				assert.Nil(t, testSerialized.Claims)
				assert.IsType(t, ErrInvalidSerializedClaims{}, err)
				assert.EqualError(t, err, "invalid serialized claims: invalid type, notCorrect")
			},
		},
		{
			name:                     "loginClaims - success",
			marshalUnmarshalTestPair: loginClaimsTestPair,
			check: func(t *testing.T, pair marshalUnmarshalTestPair, err error, testSerialized Serialized) {
				assert.Equal(
					t,
					pair.Claims,
					testSerialized.Claims,
				)
			},
		},
		{
			name: "loginClaims - unmarshall failure",
			marshalUnmarshalTestPair: marshalUnmarshalTestPair{
				Marshalled: []byte(fmt.Sprintf(
					"{\"type\":\"%s\", \"userID\":1234}",
					LoginClaimsType,
				)),
			},
			check: func(t *testing.T, pair marshalUnmarshalTestPair, err error, testSerialized Serialized) {
				assert.Nil(t, testSerialized.Claims)
				assert.IsType(t, ErrUnmarshal{}, err)
				assert.EqualError(t, err, "unmarshalling error: json: cannot unmarshal number into Go struct field Login.userID of type identifier.ID")
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				testSerialized := Serialized{}
				err := testSerialized.UnmarshalJSON(tt.marshalUnmarshalTestPair.Marshalled)
				tt.check(t, tt.marshalUnmarshalTestPair, err, testSerialized)
			},
		)
	}
}
