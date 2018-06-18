package router

import (
	"github.com/ktpswjz/httpserver/example/webserver/server/controller/admin"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/router"
)

type adminController struct {
	adminConfig *admin.Config
	adminSys *admin.Sys
}

func (s *innerRouter) mapAdminApi(path types.Path, router *router.Router) {
	s.adminConfig = &admin.Config{}
	s.adminConfig.Config = s.cfg
	s.adminConfig.SetLog(s.GetLog())
	s.adminSys = &admin.Sys{}

	// 获取配置信息
	router.POST(path.Path("/config"), s.adminConfig.GetInfo, nil)

	// 系统信息
	router.POST(path.Path("/sys/host"), s.adminSys.GetHost, s.adminSys.GetHostDoc)
	router.POST(path.Path("/sys/network/interfaces"), s.adminSys.GetNetworkInterfaces, s.adminSys.GetNetworkInterfacesDoc)
	router.POST(path.Path("/sys/disk/partitions"), s.adminSys.GetDiskPartitions, s.adminSys.GetDiskPartitionsDoc)
}