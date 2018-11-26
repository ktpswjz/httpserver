package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"strings"
)

func Decode(jwt string, payload, header interface{}) ([]string, error) {
	// jwt: header.payload.signature
	values := strings.Split(jwt, ".")
	if len(values) != 3 {
		return nil, fmt.Errorf("invalid format")
	}

	if header != nil {
		err := unmarshal(values[0], header)
		if err != nil {
			return nil, fmt.Errorf("invalid hearder: %s", err.Error())
		}
	}

	if payload != nil {
		err := unmarshal(values[1], payload)
		if err != nil {
			return nil, fmt.Errorf("invalid payload: %s", err.Error())
		}
	}

	return values, nil
}

func Encode(secret, algorithm string, payload interface{}) (string, error) {
	header := &Header{
		Algorithm: algorithm,
		Type:      "JWT",
	}

	headerBase64, err := marshal(header)
	if err != nil {
		return "", err
	}
	payloadBase64, err := marshal(payload)
	if err != nil {
		return "", err
	}

	msg := fmt.Sprintf("%s.%s", headerBase64, payloadBase64)
	signature, err := sign(secret, algorithm, msg)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s", msg, signature), nil
}

func Verify(jwt, secret string) error {
	header := &Header{}
	values, err := Decode(jwt, nil, header)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("%s.%s", values[0], values[1])
	signature, err := sign(secret, header.Algorithm, msg)
	if err != nil {
		return fmt.Errorf("sign fail: %s", err.Error())
	}
	if signature != values[2] {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

func sign(secret, algorithm, msg string) (string, error) {
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
	h.Write([]byte(msg))
	bytes := h.Sum(nil)

	return toBase64(bytes), nil
}

func toBase64(src []byte) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(src)
}

func unmarshal(base64Val string, val interface{}) error {
	data, err := base64.URLEncoding.DecodeString(base64Val)
	if err != nil {
		data, err = base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(base64Val)
		if err != nil {
			return err
		}
	}

	return json.Unmarshal(data, val)
}

func marshal(val interface{}) (string, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	return toBase64(data), nil
}
