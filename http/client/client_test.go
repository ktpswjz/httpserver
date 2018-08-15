package client

import "testing"

func TestClient_PostJson(t *testing.T) {
	url := "http://172.16.99.181/auth/rsa/key/public"
	argument := &inputArgument{
		ID: 11,
		Name: "Json",
	}
	client := &Client{}
	input, output, _, err := client.PostJson(url, argument)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("input: ", string(input[:]))
	t.Log("output:", string(output[:]))
}

type inputArgument struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}