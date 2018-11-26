package admin

import (
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/performance/disk"
	"github.com/ktpswjz/httpserver/performance/host"
	"github.com/ktpswjz/httpserver/performance/network"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"net/http"
	"time"
)

type Sys struct {
}

func (s *Sys) GetHost(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	data, err := host.Info()
	if err != nil {
		a.Error(errors.Exception, err)
		return
	}

	a.Success(data)
}

func (s *Sys) GetHostDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取主机信息")
	function.SetNote("获取服务系统当前相关信息")
	function.SetOutputExample(&host.Host{
		ID:              "8f438ea2-c26b-401e-9f6b-19f2a0e4ee2e",
		Name:            "pc",
		BootTime:        types.Time(time.Now()),
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "18.04",
		KernelVersion:   "4.15.0-22-generic",
		CPU:             "Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz x2",
		Memory:          "4GB",
	})
	function.SetContentType("")

	catalog := a.CreateCatalog("平台管理", "平台管理服务相关接口")
	catalog.CreateChild("系统信息", "系统配置相关接口").SetFunction(function)

	return function
}

func (s *Sys) GetNetworkInterfaces(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	data, err := network.Interfaces()
	if err != nil {
		a.Error(errors.Exception, err)
		return
	}

	a.Success(data)
}

func (s *Sys) GetNetworkInterfacesDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取网卡信息")
	function.SetNote("获取主机网卡相关信息")
	function.SetOutputExample([]network.Interface{
		{
			Name:    "本地连接",
			MTU:     1500,
			MacAddr: "00:16:5d:13:b9:70",
			IPAddrs: []string{
				"fe80::b1d0:ff08:1f6f:3e0b/64",
				"192.168.1.1/24",
			},
			Flags: []string{
				"up",
				"broadcast",
				"multicast",
			},
		},
	})
	function.SetContentType("")

	catalog := a.CreateCatalog("平台管理", "平台管理服务相关接口")
	catalog.CreateChild("系统信息", "系统配置相关接口").SetFunction(function)

	return function
}

func (s *Sys) GetDiskPartitions(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	data, err := disk.Partitions()
	if err != nil {
		a.Error(errors.Exception, err)
		return
	}

	a.Success(data)
}

func (s *Sys) GetDiskPartitionsDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取磁盘分区信息")
	function.SetNote("获取主机磁盘分区相关信息")
	function.SetOutputExample([]disk.Partition{
		{
			FileSystem:     "C:",
			FileSystemType: "NTFS",
			Path:           "C:",
			Total:          "120GB",
			Free:           "60.5GB",
			Used:           "59.5GB",
			UsedPercent:    49.8,
		},
	})
	function.SetContentType("")

	catalog := a.CreateCatalog("平台管理", "平台管理服务相关接口")
	catalog.CreateChild("系统信息", "系统配置相关接口").SetFunction(function)

	return function
}
