package router

import (
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/router"
	"net/http"
)

func (s *innerRouter) mapAdminSite(path types.Path, router *router.Router, root string) {
	router.ServeFiles(path.Path("/*filepath"), http.Dir(root), nil)
}