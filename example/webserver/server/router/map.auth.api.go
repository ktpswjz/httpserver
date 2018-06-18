package router

import (
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/example/webserver/server/controller/auth"
)

type authController struct {
	authAdmin *auth.Admin
}

func (s *innerRouter) mapAuthApi(path types.Path, router *router.Router) {
	s.authAdmin = &auth.Admin{}
	s.authAdmin.SetLog(s.GetLog())
	s.authAdmin.Config = s.cfg
	s.authAdmin.DbToken = s.dbToken

	// 获取平台信息
	router.POST(path.Path("/sys/info"), s.authAdmin.GetInfo, s.authAdmin.GetInfoDoc)
}
