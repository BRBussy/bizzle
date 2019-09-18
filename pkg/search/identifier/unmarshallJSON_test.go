package identifier

import (
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
				"json unmarshal",
				"invalid character 'o' in literal null (expecting 'u')",
			},
		}.Error(),
		"error should be correct for incorrect input",
	)
}

func TestUnmarshalSerializedIdentifier_InvalidType(t *testing.T) {
	assert := testifyAssert.New(t)
	assert.EqualError(
		(&Serialized{}).UnmarshalJSON([]byte("{\"id\":\"1234\"}")),
		ErrInvalidSerializedIdentifier{Reasons: []string{"no type field"}}.Error(),
		"error should be correct for input without type field",
	)

	assert.EqualError(
		(&Serialized{}).UnmarshalJSON([]byte("{\"type\":\"notAValidType\",\"id\":\"1234\"}")),
		ErrInvalidSerializedIdentifier{
			Reasons: []string{
				"invalid type",
				"\"notAValidType\"",
			},
		}.Error(),
		"error should be correct for nil input",
	)
}
