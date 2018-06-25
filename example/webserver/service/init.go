package main

import (
	"github.com/ktpswjz/httpserver/example/webserver/server/config"
	"github.com/ktpswjz/httpserver/logger"
	"github.com/ktpswjz/httpserver/types"
	"os"
	"fmt"
	"github.com/kardianos/service"
	"path/filepath"
)

const (
	moduleType    = "server"
	moduleName    = "web.server"
	moduleRemark  = "WEB服务器"
	moduleVersion = "1.0.1.0"
)

var (
	cfg = config.NewConfig()
	log = &logger.Writer{Level: logger.LevelAll}
	pro = &Program{}
	svc service.Service = nil
)

func init() {
	moduleArgs := &types.Args{}
	serverArgs := &Args{}
	moduleArgs.Parse(os.Args, moduleType, moduleName, moduleVersion, moduleRemark, serverArgs)

	// service
	svcCfg := &service.Config{
		Name: moduleName,
		DisplayName: moduleName,
		Description: moduleRemark,
	}
	configPath := serverArgs.config
	if configPath != "" {
		svcCfg.Arguments = []string{fmt.Sprintf("-config=%s", configPath)}
	}
	svcVal, err := service.New(pro, svcCfg)
	if err != nil {
		fmt.Print("init service fail: ", err)
		os.Exit(4)
	}
	svc = svcVal
	if serverArgs.help {
		serverArgs.ShowHelp()
		os.Exit(0)
	} else if serverArgs.isInstall {
		err = svc.Install()
		if err != nil {
			fmt.Println("install service ", svc.String(), " fail: ", err)
		} else {
			fmt.Println("install service ", svc.String(), " success")
		}
		os.Exit(0)
	} else if serverArgs.isUninstall {
		err = svc.Uninstall()
		if err != nil {
			fmt.Println("uninstall service ", svc.String(), " fail: ", err)
		} else {
			fmt.Println("uninstall service ", svc.String(), " success")
		}
		os.Exit(0)
	} else if serverArgs.isStart {
		err = svc.Start()
		if err != nil {
			fmt.Println("start service ", svc.String(), " fail: ", err)
		} else {
			fmt.Println("start service ", svc.String(), " success")
		}
		os.Exit(0)
	} else if serverArgs.isStop {
		err = svc.Stop()
		if err != nil {
			fmt.Println("stop service ", svc.String(), " fail: ", err)
		} else {
			fmt.Println("stop service ", svc.String(), " success")
		}
		os.Exit(0)
	} else if serverArgs.isRestart {
		err = svc.Restart()
		if err != nil {
			fmt.Println("restart service ", svc.String(), " fail: ", err)
		} else {
			fmt.Println("restart service ", svc.String(), " success")
		}
		os.Exit(0)
	}

	// config
	if configPath == "" {
		configPath = filepath.Join(moduleArgs.ModuleFolder(), "config.json")
	}
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		err = cfg.SaveToFile(configPath)
		if err != nil {
			fmt.Println("generate configure file fail: ", err)
		}
	} else {
		err = cfg.LoadFromFile(configPath)
		if err != nil {
			fmt.Println("load configure file fail: ", err)
		}
	}
	if serverArgs.showConfig {
		fmt.Println(cfg.FormatString())
		os.Exit(0)
	}
	cfg.SetArgs(moduleArgs)

	if cfg.Server.Https.Enabled {
		certFilePath := cfg.Server.Https.Cert.File
		if certFilePath == "" {
			certFilePath = filepath.Join(moduleArgs.ModuleFolder(), "server.pfx")
			cfg.Server.Https.Cert.File = certFilePath
		}
	}
	if serverArgs.showConfig {
		fmt.Println(cfg.String())
		os.Exit(0)
	}

	log.Init(cfg.Log.Level, moduleName, cfg.Log.Folder)
	log.Std = true

	//LogInfo("start at: ", moduleArgs.ModulePath())
	//LogInfo("version: ", moduleVersion)
	//LogInfo("log path: ", cfg.Log.Folder)
	//LogInfo("configure path: ", configPath)
	//LogInfo("configure info: ", cfg)
}
