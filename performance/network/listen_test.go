package network

import (
	"testing"
)

func TestListeningPorts(t *testing.T) {
	ports := ListeningPorts()
	count := len(ports)
	t.Log("count = ", count)

	for i := 0; i < count; i++ {
		item := ports[i]
		t.Logf("%3d %18s:%-6d %s", i+1, item.Address, item.Port, item.Protocol)
	}
}
