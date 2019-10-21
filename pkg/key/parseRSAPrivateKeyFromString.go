package key

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ParseRSAPrivateKeyFromString(rsaPrivateKeyString string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(rsaPrivateKeyString))
	if block == nil {
		return nil, errors.New("failed to parse private key string")
	}

	pvtKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pvtKey, nil
}
