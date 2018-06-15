package memory

import "testing"

func TestInfo(t *testing.T)  {
	v, err := Info()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("total:", v.TotalText)
	t.Log("available:", v.AvailableText)
	t.Log("used:", v.UsedText)
	t.Log("used percent:", v.UsedPercent, "%")
}
