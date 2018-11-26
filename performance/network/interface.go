package network

import "github.com/shirou/gopsutil/net"

type Interface struct {
	Name    string   `json:"name" note:"网卡名称"`
	MTU     int      `json:"mtu" note:"最大传输单元"`
	MacAddr string   `json:"macAddr" note:"MAC地址"`
	IPAddrs []string `json:"ipAddrs" note:"IP地址"`
	Flags   []string `json:"flags" note:"标志, 如up, loopback, multicast"`
}

func Interfaces() ([]Interface, error) {
	vs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	infos := make([]Interface, 0)
	for _, v := range vs {
		info := Interface{
			Name:    v.Name,
			MTU:     v.MTU,
			MacAddr: v.HardwareAddr,
			IPAddrs: make([]string, 0),
			Flags:   make([]string, 0),
		}
		ipCount := len(v.Addrs)
		for i := 0; i < ipCount; i++ {
			info.IPAddrs = append(info.IPAddrs, v.Addrs[i].Addr)
		}
		flagCount := len(v.Flags)
		for i := 0; i < flagCount; i++ {
			info.Flags = append(info.Flags, v.Flags[i])
		}

		infos = append(infos, info)
	}

	return infos, nil
}
