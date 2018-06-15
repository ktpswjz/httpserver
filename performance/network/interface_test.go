package network

import "testing"

func TestInterfaces(t *testing.T) {
	vs, err := Interfaces()
	if err != nil {
		t.Fatal(err)
	}

	for i, v := range vs {
		t.Log(i, ":", v)
	}
}
