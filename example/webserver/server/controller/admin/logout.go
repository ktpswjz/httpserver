package admin

import (
	"github.com/ktpswjz/httpserver/example/webserver/server/controller"
	"net/http"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/example/webserver/server/config"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/example/webserver/database/memory"
)

type Logout struct {
	controller.Base
}

func NewLogout(cfg *config.Config, log types.Log, dbToken memory.Token) *Logout  {
	instance := &Logout{}
	instance.Config = cfg
	instance.SetLog(log)
	instance.DbToken = dbToken

	return instance
}

func (s *Logout) Logout(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	token, err := s.DbToken.Get(a.Token())
	if err != nil {
		a.Error(errors.Exception,  err)
		return
	}
	if token == nil {
		a.Error(errors.NotExist,  "凭证'", a.Token(), "'不存在")
		return
	}
	err = s.DbToken.Del(token)
	if err != nil {
		a.Error(errors.Exception,  err)
		return
	}

	a.Success(true)
}

func (s *Logout) LogoutDoc(a document.Assistant) document.Function  {
	function := a.CreateFunction("退出登录")
	function.SetNote("退出登录, 使当前凭证失效")
	function.SetOutputExample(true)

	catalog := a.CreateCatalog("平台管理", "平台管理服务相关接口")
	catalog.CreateChild("权限管理", "系统授权相关接口").SetFunction(function)

	return function
}