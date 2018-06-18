package admin

import (
	"net/http"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/example/webserver/server/controller"
)

type Config struct {
	controller.Base
}

func (s *Config) GetInfo(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant)  {
	a.Success(s.Config)
}