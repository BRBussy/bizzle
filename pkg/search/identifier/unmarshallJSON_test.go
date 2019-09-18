package identifier

import (
	"fmt"
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalSerializedIdentifier_InvalidInput(t *testing.T) {
	assert := testifyAssert.New(t)

	assert.EqualError(
		(&Serialized{}).UnmarshalJSON(nil),
		ErrInvalidSerializedIdentifier{Reasons: []string{"json identifier data is nil"}}.Error(),
		"error should be correct for nil input",
	)

	assert.EqualError(
		(&Serialized{}).UnmarshalJSON([]byte("notValidJSON")),
		ErrUnmarshal{
			Reasons: []string{
				"json unmarshal id object",
				"invalid character 'o' in literal null (expecting 'u')",
			},
		}.Error(),
		"error should be correct for incorrect input",
	)
}

func TestUnmarshalSerializedIdentifier_InvalidType(t *testing.T) {
	assert := testifyAssert.New(t)

	// missing type field
	assert.EqualError(
		(&Serialized{}).UnmarshalJSON([]byte("{\"id\":\"1234\"}")),
		ErrInvalidSerializedIdentifier{
			Reasons: []string{"no type field"},
		}.Error(),
		"error should be correct for input without type field",
	)

	// invalid type field
	assert.EqualError(
		(&Serialized{}).UnmarshalJSON([]byte("{\"type\":1234,\"id\":\"1234\"}")),
		ErrUnmarshal{
			Reasons: []string{
				"json unmarshal id type",
				"json: cannot unmarshal number into Go value of type identifier.Type",
			},
		}.Error(),
		"error should be correct invalid type field",
	)
	assert.EqualError(
		(&Serialized{}).UnmarshalJSON([]byte("{\"type\":\"notAValidType\",\"id\":\"1234\"}")),
		ErrInvalidSerializedIdentifier{
			Reasons: []string{
				"invalid type",
				"notAValidType",
			},
		}.Error(),
		"error should be correct invalid type field",
	)
}

func TestUnmarshalSerializedIdentifier_IDIdentifierErrors(t *testing.T) {
	assert := testifyAssert.New(t)

	// invalid id identifier value
	assert.EqualError(
		(&Serialized{}).UnmarshalJSON([]byte(fmt.Sprintf(
			"{\"type\":\"%s\",\"id\":1234}",
			IDIdentifierType,
		))),
		ErrUnmarshal{
			Reasons: []string{
				"json: cannot unmarshal object into Go value of type identifier.ID",
			},
		}.Error(),
		"error should be correct for invalid value types",
	)
}
