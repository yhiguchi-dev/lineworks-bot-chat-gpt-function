package lineworks

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

type Signer interface {
	Sign(message string) ([]byte, error)
}

type signer struct {
	key *rsa.PrivateKey
}

func (s signer) Sign(message string) ([]byte, error) {
	hash := crypto.Hash.New(crypto.SHA256)
	hash.Write([]byte(message))
	hashed := hash.Sum(nil)

	signed, err := rsa.SignPKCS1v15(rand.Reader, s.key, crypto.SHA256, hashed)
	if err != nil {
		return nil, err
	}
	return signed, nil
}

func NewSigner(key *rsa.PrivateKey) Signer {
	return &signer{key}
}
