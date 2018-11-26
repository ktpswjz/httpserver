package rsakey

import "testing"

func TestSignData(t *testing.T) {
	pk := &Private{}
	err := pk.Create(2048)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte(string("testsss"))
	signedData, err := pk.SignData(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("sign:", signedData)
	t.Log("sign length:", len(signedData))
}
