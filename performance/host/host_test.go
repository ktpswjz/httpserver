package host

import "testing"

func TestInfo(t *testing.T)  {
	v, err := Info()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("id:", v.ID)
	t.Log("name:", v.Name)
	t.Log("boot time:", v.BootTime)
	t.Log("os:", v.OS)
	t.Log("platform:", v.Platform)
	t.Log("platform version:", v.PlatformVersion)
	t.Log("kernel version:", v.KernelVersion)
	t.Log("cpu:", v.CPU)
	t.Log("mem:", v.Memory)
}
