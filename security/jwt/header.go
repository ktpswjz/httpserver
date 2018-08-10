package jwt

import (
	"encoding/json"
	"encoding/base64"
)

type Header struct {
	Algorithm	string	`json:"alg"`
	Type		string	`json:"typ"`
}

func (s Header) String() string  {
	bytes, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes)
}
