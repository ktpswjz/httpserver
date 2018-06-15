package disk

import (
	"github.com/shirou/gopsutil/disk"
	"fmt"
)

type Partition struct {
	FileSystem	 	string	`json:"fileSystem" note:"文件系统, 如/dev/sda1"`
	FileSystemType 	string	`json:"fileSystemType" note:"文件系统类型, 如NTFS"`
	Path			string	`json:"path" note:"路径(挂载点), 如/sys"`
	Total           string  `json:"total" note:"总容量"`
	Free            string  `json:"free" note:"剩余容量"`
	Used            string  `json:"used" note:"已用容量"`
	UsedPercent     float64 `json:"usedPercent" note:"已用比分比"`
}

func Partitions() ([]Partition, error)  {
	ps, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	infos := make([]Partition, 0)
	for _, p := range ps {
		u, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}

		info := Partition{Path: u.Path}
		info.FileSystem = p.Device
		info.FileSystemType = p.Fstype
		info.Total = diskToText(float64(u.Total))
		info.Free = diskToText(float64(u.Free))
		info.Used = diskToText(float64(u.Used))
		info.UsedPercent = u.UsedPercent

		infos = append(infos, info)
	}



	return infos, nil
}


func diskToText(v float64) string  {
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