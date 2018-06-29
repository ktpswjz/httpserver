package router

import (
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/example/webserver/server/config"
	"github.com/ktpswjz/httpserver/router"
	"net/http"
	"strings"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/database/memory"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"time"
)

const (
	App  = "/app"
	Auth    = "/auth"
	Doc  = "/doc"
	DocApi  = "/doc.api"
	Admin  = "/admin"
	AdminApi  = "/api"
	Favicon = "/admin/favicon.ico"
)

type Router interface {
	Map(router *router.Router)
	PreRouting(w http.ResponseWriter, r *http.Request, a router.Assistant) bool
	PostRouting(a router.Assistant)
}

func NewRouter(cfg *config.Config, log types.Log) Router  {
	instance := &innerRouter {cfg: cfg}
	instance.SetLog(log)
	instance.dbToken, _ = memory.NewToken(cfg.Site.Admin.Api.Token.Expiration, log)

	return instance
}

type innerRouter struct {
	types.Base
	cfg *config.Config
	dbToken memory.Token

	// controllers
	docController
	authController
	adminController
}

func (s *innerRouter) Map(router *router.Router) {
	router.Doc = document.NewDocument(s.cfg.Site.Doc.Enable, s.GetLog())

	s.mapAuthApi(types.Path{Prefix:Auth}, router)
	s.mapAppSite(types.Path{Prefix:App}, router, s.cfg.Site.App.Root)

	if s.cfg.Site.Admin.Enable {
		s.mapAdminApi(types.Path{Prefix:AdminApi}, router)
		s.mapAdminSite(types.Path{Prefix:Admin}, router, s.cfg.Site.Admin.Root)
		s.LogInfo("admin is enabled")
	}

	if s.cfg.Site.Doc.Enable {
		s.mapDocApi(types.Path{Prefix:DocApi}, router)
		s.mapDocSite(types.Path{Prefix:Doc}, router, s.cfg.Site.Doc.Root)
		s.LogInfo("doc is enabled")
		router.Doc.GenerateCatalogTree()
	}
}

func (s *innerRouter) PreRouting(w http.ResponseWriter, r *http.Request, a router.Assistant) bool {
	path := r.URL.Path
	if isApi(path) {
		// enable across access
		if r.Method == "OPTIONS" {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "content-type,token")
			return true
		}
	}

	if isAdminApi(path) {
		token := a.Token()
		if token == "" {
			a.Error(errors.AuthNoToken)
			return true
		}
		tokenEntity, err := s.dbToken.Get(token)
		if err != nil || tokenEntity == nil {
			a.Error(errors.AuthTokenInvalid, err)
			return true
		}
		if a.RIP() != tokenEntity.LoginIP {
			a.Error(errors.AuthTokenIllegal)
			return true
		}
		tokenEntity.ActiveTime = time.Now()
	}

	// default to admin site
	if "/" == r.URL.Path {
		r.URL.Path = Admin
	} else if "/favicon.ico" == r.URL.Path {
		r.URL.Path = Favicon
	}

	return false
}

func (s *innerRouter) PostRouting(a router.Assistant) {

}


func isApi(path string) bool  {
	return strings.HasPrefix(path, Auth) || strings.HasPrefix(path, AdminApi)
}

func isAdminApi(path string) bool  {
	return strings.HasPrefix(path, AdminApi)
}
