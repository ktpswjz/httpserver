package admin

import (
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/database/memory"
	"github.com/ktpswjz/httpserver/example/webserver/model"
	"github.com/ktpswjz/httpserver/example/webserver/server/config"
	"github.com/ktpswjz/httpserver/example/webserver/server/controller"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"net/http"
	"time"
)

type Service struct {
	controller.Base

	bootTime time.Time
}

func NewService(cfg *config.Config, log types.Log, dbToken memory.Token) *Service {
	instance := &Service{}
	instance.Config = cfg
	instance.SetLog(log)
	instance.DbToken = dbToken
	instance.bootTime = time.Now()

	return instance
}

func (s *Service) GetInfo(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	data := &model.ServiceInfo{BootTime: types.Time(s.bootTime)}
	if s.Config != nil {
		args := s.Config.GetArgs()
		if args != nil {
			data.Name = args.ModuleName()
			data.Version = args.ModuleVersion().ToString()
			data.Remark = args.ModuleRemark()
		}
	}

	a.Success(data)
}

func (s *Service) GetInfoDoc(a document.Assistant) document.Function {
	function := a.CreateFunction("获取当前服务信息")
	function.SetNote("获取当前服务信息")
	function.SetOutputExample(&model.ServiceInfo{
		Name:     "server",
		BootTime: types.Time(time.Now()),
		Version:  "1.0.1.0",
		Remark:   "XXX服务",
	})
	function.SetContentType("")

	catalog := a.CreateCatalog("平台管理", "平台管理服务相关接口")
	catalog.CreateChild("服务信息", "服务当前相关接口").SetFunction(function)

	return function
}
