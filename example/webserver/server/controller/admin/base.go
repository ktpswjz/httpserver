package admin

import (
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/example/webserver/server/config"
)

type Base struct {
	types.Base

	Config *config.Config
}