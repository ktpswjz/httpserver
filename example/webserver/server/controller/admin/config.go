package admin

import (
	"net/http"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/example/webserver/server/controller"
	"github.com/ktpswjz/httpserver/example/webserver/server/config"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/example/webserver/database/memory"
)

type Config struct {
	controller.Base
}

func NewConfig(cfg *config.Config, log types.Log, dbToken memory.Token) *Config  {
	instance := &Config{}
	instance.Config = cfg
	instance.SetLog(log)
	instance.DbToken = dbToken

	return instance
}

func (s *Config) GetInfo(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant)  {
	a.Success(s.Config)
}