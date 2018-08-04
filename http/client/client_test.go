package client

import "testing"

func TestClient_PostJson(t *testing.T) {
	url := "http://172.16.99.181/auth/rsa/key/public"
	client := &Client{}
	bytes, _, err := client.PostJson(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bytes[:]))
}
