package jwt

import (
	"fmt"
)

type Encoding interface {
	Encode(payload interface{}) (string, error)
	Decode(value string, payload interface{}) error
}

func NewEncoding(algorithm, secret string) Encoding  {
	return &encoding{
		algorithm: algorithm,
		secret: secret,
	}
}

type encoding struct {
	algorithm string
	secret    string
}

func (s *encoding) Decode(jwt string, payload interface{}) error {
	// jwt: header.payload.signature
	header := &Header{}
	values, err := Decode(jwt, payload, header)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("%s.%s", values[0], values[1])
	signature, err := sign(s.secret, s.algorithm, msg)
	if err != nil {
		return fmt.Errorf("sign fail: %s", err.Error())
	}
	if signature != values[2] {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

func (s *encoding) Encode(payload interface{}) (string, error) {
	return Encode(s.secret, s.algorithm, payload)
}

