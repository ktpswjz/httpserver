package disk

import "testing"

func TestPartitions(t *testing.T) {
	vs, err := Partitions()
	if err != nil {
		t.Fatal(err)
	}

	for i, v := range vs {
		t.Log(i, ":", v)
	}
}
