package claims

import (
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLogin_Type(t *testing.T) {
	assert := testifyAssert.New(t)

	assert.Equal(
		LoginClaimsType,
		Login{}.Type(),
	)
}

func TestLogin_ToJSON(t *testing.T) {
	assert := testifyAssert.New(t)

	jsonClaims, err := loginClaimsTestPair.Claims.ToJSON()

	assert.Nil(err)
	assert.JSONEq(
		string(loginClaimsTestPair.Marshalled),
		string(jsonClaims),
	)
}

func TestLogin_Expired(t *testing.T) {
	assert := testifyAssert.New(t)
	typedClaims, ok := loginClaimsTestPair.Claims.(Login)
	assert.True(ok)

	// set expiration in the past
	typedClaims.ExpirationTime = time.Date(
		time.Now().Year()-1,
		time.Now().Month(),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		time.Now().Nanosecond(),
		time.Now().Location(),
	).UTC().Unix()
	assert.Equal(
		true,
		typedClaims.Expired(),
	)

	// set expiration in future
	typedClaims.ExpirationTime = time.Now().Add(time.Hour * 2).UTC().Unix()
	assert.Equal(
		false,
		typedClaims.Expired(),
	)
}

func TestLogin_ExpiryTime(t *testing.T) {
	assert := testifyAssert.New(t)
	typedClaims, ok := loginClaimsTestPair.Claims.(Login)
	assert.True(ok)

	assert.Equal(
		typedClaims.ExpirationTime,
		typedClaims.ExpiryTime(),
	)
}
