package jwt

import "testing"

func TestJwt_Encode(t *testing.T) {
	jwt := NewEncoding("HS256", "your-256-bit-secret")

	payload := &Payload{
		Sub:  "1234567890",
		Name: "John Doe",
		Iat:  1516239022,
	}
	sign := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	signature, err := jwt.Encode(payload)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("expect: ", sign)
	t.Log("actual: ", signature)

	if signature != sign {
		t.Error("Invalid Signature")
	}
}

func TestJwt_Decode(t *testing.T) {
	jwt := NewEncoding("HS384", "your-384-bit-secret")
	//payload := &Payload{
	//	Sub: "11",
	//	Name: "22",
	//	Iat: 44,
	//}
	payload := &Payload{}
	value := "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMSIsIm5hbWUiOiIyMiIsImlhdCI6NDR9.2HRRmEiyxZcmchZpBNeqdaVaHkncfke1FrRV7r9AM0Y-QPWB9IfpkfOsQY7Ou8vu"

	err := jwt.Decode(value, payload)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("expect: ", "11")
	t.Log("actual: ", payload.Sub)

	if payload.Sub != "11" {
		t.Error("Invalid Signature")
	}
	if payload.Name != "22" {
		t.Error("Invalid Signature")
	}
	if payload.Iat != 44 {
		t.Error("Invalid Signature")
	}

	value = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJWSUNPLkhCUC5TRVJWSUNFIiwic3ViIjoiVklDTy5IQlAuU0VSVklDRSIsImF1ZCI6Im9yZ2FuSWQiLCJleHAiOjE1MzQ4MzIwMzYyMzYsIm5iZiI6MTUzNDgzMTQzNjIzNiwiaWF0IjoxNTM0ODMxNDM2MjM2LCJqdGkiOiI5MWMzN2IxMi1lNmZkLTQ2ODItYjY4NC1iMDFiOTE4MDkwMTIifQ.em_1Nvo6J5MNljBAZHf6nDSL-hU56_aS_icjqLyvvCw"
	_, err = Decode(value, payload, nil)
	if err != nil {
		t.Fatal(err)
	}
}

type Payload struct {
	Sub  string `json:"sub"`
	Name string `json:"name"`
	Iat  int    `json:"iat"`
}
