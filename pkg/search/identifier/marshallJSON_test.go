package identifier

import (
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalSerializedIdentifier_MarshalSuccesses(t *testing.T) {
	assert := testifyAssert.New(t)

	for _, tc := range allTestCases {
		testSerialized := Serialized{
			Identifier: tc.Identifier,
		}
		marshalledJSON, err := testSerialized.MarshalJSON()
		assert.Equal(
			nil,
			err,
			"error should be nil when marshalling",
		)
		assert.JSONEq(
			string(tc.SerializedIdentifier),
			string(marshalledJSON),
			"marshalled json should match expected",
		)
	}
}
