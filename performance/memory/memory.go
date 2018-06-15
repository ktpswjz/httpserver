package memory

import (
	"github.com/shirou/gopsutil/mem"
	"fmt"
)

type Memory struct {
	Total uint64 `json:"total"`
	TotalText string `json:"totalText"`
	Available uint64 `json:"available"`
	AvailableText string `json:"availableText"`
	Used uint64 `json:"used"`
	UsedText string `json:"usedText"`
	UsedPercent float64 `json:"usedPercent"`
}

func Info()(*Memory, error)  {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	info := &Memory{
		Total: v.Total,
		Available: v.Available,
		Used: v.Used,
		UsedPercent: v.UsedPercent,
	}
	info.TotalText = memToText(float64(v.Total))
	info.AvailableText = memToText(float64(v.Available))
	info.UsedText = memToText(float64(v.Used))

	return info, nil
}


func memToText(v float64) string  {
	kb := float64(1024)
	mb := 1024 * kb
	gb := 1024 * mb

	if v >= gb {
		return fmt.Sprintf("%.1fGB", v / gb)
	} else if v >= mb {
		return fmt.Sprintf("%.1fMB", v / mb)
	} else if v >= kb {
		return fmt.Sprintf("%.1fKB", v / kb)
	} else {
		return fmt.Sprintf("%.0fB", v)
	}
}