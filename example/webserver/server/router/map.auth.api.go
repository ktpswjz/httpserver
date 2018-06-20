package router

import (
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/example/webserver/server/controller/auth"
	"github.com/ktpswjz/httpserver/example/webserver/server/authentication"
)

type authController struct {
	authAdmin *auth.Admin
	adminAuthentication authentication.Admin
}

func (s *innerRouter) mapAuthApi(path types.Path, router *router.Router) {
	s.adminAuthentication.Config = &s.cfg.Site.Admin
	s.authAdmin = &auth.Admin{}
	s.authAdmin.SetLog(s.GetLog())
	s.authAdmin.Config = s.cfg
	s.authAdmin.DbToken = s.dbToken
	s.authAdmin.Authenticate = s.adminAuthentication.Authenticate

	// 获取平台信息
	router.POST(path.Path("/sys/info"), s.authAdmin.GetInfo, s.authAdmin.GetInfoDoc)

	// 获取验证码
	router.POST(path.Path("/captcha"), s.authAdmin.GetCaptcha, s.authAdmin.GetCaptchaDoc)

	// 服务管理登录
	if s.cfg.Site.Admin.Enable {
		router.POST(path.Path("/login"), s.authAdmin.Login, s.authAdmin.LoginDoc)
	}
}
