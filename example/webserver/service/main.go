package main

import (
	"fmt"
	"github.com/kardianos/service"
	"github.com/ktpswjz/httpserver/example/webserver/server/router"
	"github.com/ktpswjz/httpserver/http/server/host"
)

func main() {
	log.Std = false
	defer log.Close()

	LogInfo("start at: ", cfg.GetArgs().ModulePath())
	LogInfo("version: ", moduleVersion)
	LogInfo("log path: ", cfg.Log.Folder)
	LogInfo("configure info: ", cfg)

	router := router.NewRouter(cfg, log)
	if service.Interactive() {
		host := host.NewHost(cfg.GetServer(), router, log, nil)
		err := host.Run()
		if err != nil {
			fmt.Println("run server error:", err)
		}
	} else {
		pro.server = host.NewHost(cfg.GetServer(), router, log, svc.Restart)
		err := svc.Run()
		if err != nil {
			LogError("run service error:", err)
		}
	}
}
