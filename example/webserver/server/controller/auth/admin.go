package auth

import (
	"net/http"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/model"
	"github.com/ktpswjz/httpserver/example/webserver/server/controller"
)

type Admin struct {
	controller.Base
}

func (s *Admin) GetInfo(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant)  {
	data := &model.SysInfo {
		Name: s.Config.Name,
		BackVersion: s.Config.GetVersion().ToString(),
	}
	a.Success(data)
}

func (s *Admin) GetInfoDoc(a document.Assistant) document.Function  {
	function := a.CreateFunction("获取系统信息")
	function.SetNote("获取服务系统当前相关信息")
	function.SetOutputExample(&model.SysInfo {
		Name: "WEB服务器",
		BackVersion: "1.0.1.0",
		FrontVersion: "1.0.1.8",
	})
	function.IgnoreToken(true)
	function.SetContentType("")

	catalog := a.CreateCatalog("平台信息", "平台信息相关接口")
	catalog.SetFunction(function)

	return function
}