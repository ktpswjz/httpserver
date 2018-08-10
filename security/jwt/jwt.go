package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"strings"
	"hash"
	"crypto/sha512"
	"fmt"
	"encoding/base64"
	"encoding/json"
)

type JWT interface {
	Encode(payload interface{}) (string, error)
	Decode(value string, payload interface{}) error
}

func New(algorithm, secret string) JWT  {
	return &jwt{
		algorithm: algorithm,
		secret: secret,
	}
}

type jwt struct {
	algorithm	string
	secret		string
}

func (s *jwt) Decode(value string, payload interface{}) error {
	// header.payload.signature
	values := strings.Split(value, ".")
	if len(values) != 3 {
		return fmt.Errorf("invalid format")
	}

	headerData, err := base64.URLEncoding.DecodeString(values[0])
	if err != nil {
		return fmt.Errorf("invalid format, decode header fail: %s", err.Error())
	}
	header := &Header{}
	err = json.Unmarshal(headerData, header)
	if err != nil {
		return fmt.Errorf("invalid format, unmarshal header fail: %s", err.Error())
	}
	headerAndPayload := fmt.Sprintf("%s.%s", values[0], values[1])
	signature, err := s.sign(s.secret, s.algorithm, headerAndPayload)
	if err != nil {
		return fmt.Errorf("sign fail: %s", err.Error())
	}
	if signature != values[2] {
		return fmt.Errorf("invalid signature")
	}

	payloadData, err := base64.URLEncoding.DecodeString(values[1])
	if err != nil {
		return fmt.Errorf("invalid format, decode payload fail: %s", err.Error())
	}
	err = json.Unmarshal(payloadData, payload)
	if err != nil {
		return fmt.Errorf("invalid format, unmarshal payloadData fail: %s", err.Error())
	}

	return nil
}

func (s *jwt) Encode(payload interface{}) (string, error)  {
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("invalid payload: %s", err.Error())
	}
	header := &Header{
		Algorithm: s.algorithm,
		Type: "JWT",
	}
	headerAndPayload := fmt.Sprintf("%s.%s", header.String(), base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(payloadData))
	signature, err := s.sign(s.secret, s.algorithm, headerAndPayload)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s", headerAndPayload, signature), nil
}

func (s *jwt) sign(secret, algorithm, headerAndPayload string) (string, error)  {
	key := []byte(secret)
	var h hash.Hash = nil
	if strings.ToUpper(algorithm) == "HS256" {
		h = hmac.New(sha256.New, key)
	} else if strings.ToUpper(algorithm) == "HS384" {
		h = hmac.New(sha512.New384, key)
	} else if strings.ToUpper(algorithm) == "HS512" {
		h = hmac.New(sha512.New, key)
	}

	if h == nil {
		return "", fmt.Errorf("algorithm '%s' not support", algorithm)
	}
	h.Write([]byte(headerAndPayload))
	bytes := h.Sum(nil)

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes), nil
}