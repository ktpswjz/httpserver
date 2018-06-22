package router

import (
	"github.com/ktpswjz/httpserver/example/webserver/server/controller/admin"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/router"
)

type adminController struct {
	adminConfig *admin.Config
	adminSys *admin.Sys
	adminLogout *admin.Logout
	adminApp *admin.App
}

func (s *innerRouter) mapAdminApi(path types.Path, router *router.Router) {
	s.adminConfig = &admin.Config{}
	s.adminConfig.Config = s.cfg
	s.adminConfig.SetLog(s.GetLog())

	s.adminSys = &admin.Sys{}

	s.adminLogout = &admin.Logout{}
	s.adminLogout.Config = s.cfg
	s.adminLogout.DbToken = s.dbToken
	s.adminLogout.SetLog(s.GetLog())

	s.adminApp = &admin.App{}
	s.adminApp.Config = s.cfg
	s.adminApp.DbToken = s.dbToken
	s.adminApp.SetLog(s.GetLog())

	// 退出登陆
	router.POST(path.Path("/logout"), s.adminLogout.Logout, s.adminLogout.LogoutDoc)

	// 获取配置信息
	router.POST(path.Path("/config"), s.adminConfig.GetInfo, nil)

	// 系统信息
	router.POST(path.Path("/sys/host"), s.adminSys.GetHost, s.adminSys.GetHostDoc)
	router.POST(path.Path("/sys/network/interfaces"), s.adminSys.GetNetworkInterfaces, s.adminSys.GetNetworkInterfacesDoc)
	router.POST(path.Path("/sys/disk/partitions"), s.adminSys.GetDiskPartitions, s.adminSys.GetDiskPartitionsDoc)

	// APP
	router.POST(path.Path("/app/upload"), s.adminApp.Upload, nil)
	router.POST(path.Path("/app/tree"), s.adminApp.Tree, s.adminApp.TreeDoc)
	router.POST(path.Path("/app/delete"), s.adminApp.Delete, s.adminApp.DeleteDoc)
}