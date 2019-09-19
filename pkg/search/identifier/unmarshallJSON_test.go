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

	// missing value
	assert.EqualError(
		(&Serialized{}).UnmarshalJSON([]byte(fmt.Sprintf(
			"{\"type\":\"%s\"}",
			IDIdentifierType,
		))),
		ErrInvalidSerializedIdentifier{
			Reasons: []string{
				"no id field",
			},
		}.Error(),
		"error should be correct for invalid value types",
	)

	// invalid id identifier value type
	assert.EqualError(
		(&Serialized{}).UnmarshalJSON([]byte(fmt.Sprintf(
			"{\"type\":\"%s\",\"id\":1234}",
			IDIdentifierType,
		))),
		ErrUnmarshal{
			Reasons: []string{
				"json: cannot unmarshal number into Go value of type identifier.ID",
			},
		}.Error(),
		"error should be correct for invalid value types",
	)

	// invalid id identifier value
	assert.EqualError(
		(&Serialized{}).UnmarshalJSON([]byte(fmt.Sprintf(
			"{\"type\":\"%s\",\"id\":\"\"}",
			IDIdentifierType,
		))),
		ErrInvalidIdentifier{
			Reasons: []string{
				"ID identifier is blank",
			},
		}.Error(),
		"error should be correct for invalid identifier",
	)
}

func TestUnmarshalSerializedIdentifier_NameIdentifierErrors(t *testing.T) {
	assert := testifyAssert.New(t)

	// missing value
	assert.EqualError(
		(&Serialized{}).UnmarshalJSON([]byte(fmt.Sprintf(
			"{\"type\":\"%s\"}",
			NameIdentifierType,
		))),
		ErrInvalidSerializedIdentifier{
			Reasons: []string{
				"no name field",
			},
		}.Error(),
		"error should be correct for invalid value types",
	)

	// invalid id identifier value type
	assert.EqualError(
		(&Serialized{}).UnmarshalJSON([]byte(fmt.Sprintf(
			"{\"type\":\"%s\",\"name\":1234}",
			NameIdentifierType,
		))),
		ErrUnmarshal{
			Reasons: []string{
				"json: cannot unmarshal number into Go value of type identifier.Name",
			},
		}.Error(),
		"error should be correct for invalid value types",
	)

	// invalid name identifier value
	assert.EqualError(
		(&Serialized{}).UnmarshalJSON([]byte(fmt.Sprintf(
			"{\"type\":\"%s\",\"name\":\"\"}",
			NameIdentifierType,
		))),
		ErrInvalidIdentifier{
			Reasons: []string{
				"Name identifier is blank",
			},
		}.Error(),
		"error should be correct for invalid identifier",
	)
}

func TestUnmarshalSerializedIdentifier_UnmarshalSuccesses(t *testing.T) {
	assert := testifyAssert.New(t)

	for _, tc := range allTestCases {
		testSerialized := Serialized{}

		assert.Equal(
			nil,
			(&testSerialized).UnmarshalJSON(tc.SerializedIdentifier),
			fmt.Sprintf("%s: error should be nil when unmarshalling", tc.id),
		)

		assert.Equal(
			tc.Identifier,
			testSerialized.Identifier,
			fmt.Sprintf("%s: identifier should equal expected", tc.id),
		)
	}
}
