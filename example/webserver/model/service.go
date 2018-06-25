package model

import "github.com/ktpswjz/httpserver/types"

type ServiceInfo struct {
	Name		string		`json:"name" note:"服务名称"`
	Version		string		`json:"version" note:"版本号"`
	BootTime	types.Time	`json:"bootTime" note:"启动时间"`
	Remark		string		`json:"remark" note:"说明"`
}
