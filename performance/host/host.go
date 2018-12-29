package host

import (
	"fmt"
	"github.com/ktpswjz/httpserver/types"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"strings"
	"time"
)

type Host struct {
	ID              string     `json:"id" note:"主机标识"`
	Name            string     `json:"name" note:"主机名称"`
	BootTime        types.Time `json:"bootTime" note:"系统启动时间"`
	OS              string     `json:"os" note:"操作系统, 如windows, linux"`
	Platform        string     `json:"platform" note:"系统平台, 如ubuntu, Microsoft Windows 10 企业版"`
	PlatformVersion string     `json:"platformVersion" note:"平台版本, 如10.0.17134 Build 17134"`
	KernelVersion   string     `json:"kernelVersion" note:"内核版本, 如4.15.0-22-generic"`
	CPU             string     `json:"cpu" note:"处理器"`
	Memory          string     `json:"memory" note:"系统内存, 如16GB"`
	TimeZone        string     `json:"timeZone" note:"系统时区, 如CST+8"`
}

func Info() (*Host, error) {
	v, err := host.Info()
	if err != nil {
		return nil, err
	}

	info := &Host{
		ID:              v.HostID,
		Name:            v.Hostname,
		BootTime:        types.Time(time.Unix(int64(v.BootTime), 0)),
		OS:              v.OS,
		Platform:        v.Platform,
		PlatformVersion: v.PlatformVersion,
		KernelVersion:   v.KernelVersion,
	}

	c, err := cpu.Info()
	if err == nil {
		if len(c) > 0 {
			info.CPU = fmt.Sprintf("%s x%d", c[0].ModelName, len(c))
		}
	}
	m, err := mem.VirtualMemory()
	if err == nil {
		info.Memory = fmt.Sprintf("%s / %s", memoryToText(float64(m.Used)), memoryToText(float64(m.Total)))
	}

	zoneName, zoneOffset := time.Now().Local().Zone()
	timeZone := strings.Builder{}
	timeZone.WriteString(zoneName)
	if zoneOffset >= 0 {
		timeZone.WriteString("+")
	}
	timeZone.WriteString(fmt.Sprint(zoneOffset / 60 / 60))
	info.TimeZone = timeZone.String()

	return info, nil
}

func memoryToText(v float64) string {
	kb := float64(1024)
	mb := 1024 * kb
	gb := 1024 * mb

	if v >= gb {
		return fmt.Sprintf("%.1fGB", v/gb)
	} else if v >= mb {
		return fmt.Sprintf("%.1fMB", v/mb)
	} else if v >= kb {
		return fmt.Sprintf("%.1fKB", v/kb)
	} else {
		return fmt.Sprintf("%.0fB", v)
	}
}
