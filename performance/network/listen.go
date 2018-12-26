package network

import (
	"sort"
	"strings"
)

type Listen struct {
	Address  string `json:"address" note:"地址"`
	Port     int    `json:"port" note:"端口"`
	Protocol string `json:"protocol" note:"协议"`
}

type ListenCollection []*Listen

func (s ListenCollection) Len() int {
	return len(s)
}
func (s ListenCollection) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ListenCollection) Less(i, j int) bool {
	if s[i].Port == s[j].Port {
		return strings.Compare(s[i].Address, s[j].Address) < 0
	}
	return s[i].Port < s[j].Port
}

func ListeningPorts() []*Listen {
	listens := make(ListenCollection, 0)
	getListenPorts(&listens)
	sort.Stable(listens)

	return listens
}
