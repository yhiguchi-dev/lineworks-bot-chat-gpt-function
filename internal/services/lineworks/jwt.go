package lineworks

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type Jwt interface {
	Create() (string, error)
}

type jwt struct {
	signer  Signer
	issuer  string
	subject string
}

func (jwt jwt) Create() (string, error) {
	header, err := createHeader()
	if err != nil {
		return "", err
	}
	claimSet, err := createClaimSet(jwt.issuer, jwt.subject)
	if err != nil {
		return "", err
	}
	signature, err := createSignature(header, claimSet, jwt.signer)
	if err != nil {
		return "", err
	}
	slice := []string{header, claimSet, signature}
	return strings.Join(slice, "."), nil
}

func NewJwt(signer Signer, issuer string, subject string) Jwt {
	return &jwt{signer, issuer, subject}
}

func createHeader() (string, error) {
	header := Header{Alg: "RS256", Typ: "JWT"}
	bytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func createClaimSet(issuer string, subject string) (string, error) {
	now := time.Now()
	expire := now.Add(time.Hour)
	claimSet := ClaimSet{Iss: issuer, Sub: subject, Iat: now.Unix(), Exp: expire.Unix()}
	bytes, err := json.Marshal(claimSet)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func createSignature(header string, claimSet string, signer Signer) (string, error) {
	slice := []string{header, claimSet}
	message := strings.Join(slice, ".")

	signed, err := signer.Sign(message)
	if err != nil {
		return "", err
	}

	signature := base64.RawURLEncoding.EncodeToString(signed)
	return signature, nil
}

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type ClaimSet struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
}
