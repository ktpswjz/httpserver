package hash

import "testing"

func TestHash(t *testing.T) {
	data := ""
	hashed, err := Hash(data, 0)
	if err != nil {
		t.Fatal(err)
	}
	if hashed != "" {
		t.Fatal(hashed)
	}

	data = ""
	hashed, err = Hash(data, MD5)
	if err != nil {
		t.Fatal(err)
	}
	if hashed != "" {
		t.Fatal(hashed)
	}

	data = " "
	hashed, err = Hash(data, MD5)
	if err != nil {
		t.Fatal(err)
	}
	if hashed != "7215ee9c7d9dc229d2921a40e899ec5f" {
		t.Fatal(hashed)
	}
	hashed, err = Hash(data, SHA1)
	if err != nil {
		t.Fatal(err)
	}
	if hashed != "b858cb282617fb0956d960215c8e84d1ccf909c6" {
		t.Fatal(hashed)
	}
	hashed, err = Hash(data, SHA256)
	if err != nil {
		t.Fatal(err)
	}
	if hashed != "36a9e7f1c95b82ffb99743e0c5c4ce95d83c9a430aac59f84ef3cbfab6145068" {
		t.Fatal(hashed)
	}
	hashed, err = Hash(data, SHA384)
	if err != nil {
		t.Fatal(err)
	}
	if hashed != "588016eb10045dd85834d67d187d6b97858f38c58c690320c4a64e0c2f92eebd9f1bd74de256e8268815905159449566" {
		t.Fatal(hashed)
	}
	hashed, err = Hash(data, SHA512)
	if err != nil {
		t.Fatal(err)
	}
	if hashed != "f90ddd77e400dfe6a3fcf479b00b1ee29e7015c5bb8cd70f5f15b4886cc339275ff553fc8a053f8ddc7324f45168cffaf81f8c3ac93996f6536eef38e5e40768" {
		t.Fatal(hashed)
	}
}
