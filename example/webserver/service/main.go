package main

import (
	"github.com/ktpswjz/httpserver/example/webserver/server/router"
	"github.com/ktpswjz/httpserver/http/server/host"
	"fmt"
	"github.com/kardianos/service"
)

func main() {
	log.Std = false
	defer log.Close()

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
