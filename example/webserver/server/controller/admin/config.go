package admin

import (
	"net/http"
	"github.com/ktpswjz/httpserver/router"
)

type Config struct {
	Base
}

func (s *Config) GetInfo(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant)  {
	a.Success(s.Config)
}