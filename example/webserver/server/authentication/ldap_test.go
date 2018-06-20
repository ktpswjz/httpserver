package authentication

import "testing"

func TestAuthenticate(t *testing.T) {
	ldap := &Ldap{
		Host: "192.168.123.1",
		Port: 389,
		Base: "dc=csby, dc=studio",
	}

	account := "test"
	password := "Dev2018"

	err := ldap.Authenticate(account, password)
	if err != nil {
		t.Fatal(err)
	}
}