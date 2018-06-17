package router

import (
	"github.com/ktpswjz/httpserver/example/webserver/server/controller/admin"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/router"
)

type adminController struct {
	adminConfig *admin.Config
}

func (s *innerRouter) mapAdminApi(path types.Path, router *router.Router) {
	s.adminConfig = &admin.Config{}
	s.adminConfig.Config = s.cfg
	s.adminConfig.SetLog(s.GetLog())

	// 获取配置信息
	router.POST(path.Path("/config"), s.adminConfig.GetInfo, nil)
}