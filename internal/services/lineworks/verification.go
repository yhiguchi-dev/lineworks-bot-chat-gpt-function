package lineworks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type Verifier interface {
	Verify(message []byte) bool
}

type verifier struct {
	key       string
	signature string
}

func (v verifier) Verify(message []byte) bool {
	mac := hmac.New(sha256.New, []byte(v.key))
	mac.Write(message)
	encoded := base64.RawStdEncoding.EncodeToString(mac.Sum(nil))
	fmt.Println(encoded)
	return string(message) == encoded
}

func NewVerifier(key string, signature string) Verifier {
	return &verifier{key, signature}
}
