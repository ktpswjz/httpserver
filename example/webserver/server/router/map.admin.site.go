package router

import (
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/types"
	"net/http"
)

func (s *innerRouter) mapAdminSite(path types.Path, router *router.Router, root string) {
	router.ServeFiles(path.Path("/*filepath"), http.Dir(root), nil)
}
