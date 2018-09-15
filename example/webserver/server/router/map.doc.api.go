package router

import (
	"github.com/ktpswjz/httpserver/example/webserver/server/controller/doc"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
)

type docController struct {
	adminDoc *doc.Doc
}

func (s *innerRouter) mapDocApi(path types.Path, router *router.Router) {
	s.adminDoc = &doc.Doc{Document: router.Doc}
	s.adminDoc.SetLog(s.GetLog())
	s.adminDoc.Config = s.cfg
	s.adminDoc.DbToken = s.dbToken

	// 接口测试
	router.POST(path.Path("/test/api"), s.adminDoc.TestApi, nil)
	router.GET(path.Path("/test/api"), s.adminDoc.TestApi, nil)

	// 获取接口目录信息
	router.POST(path.Path("/catalog/tree"), s.adminDoc.GetCatalogTree, nil)

	// 获取接口定义信息
	router.POST(path.Path("/function/:id"), s.adminDoc.GetFunction, nil)

	// 获取用户凭证
	router.POST(path.Path("/token"), s.adminDoc.CreateToken, nil)
}
