package doc

import (
	"github.com/ktpswjz/httpserver/document"
	"net/http"
	"github.com/ktpswjz/httpserver/router"
	"fmt"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/example/webserver/model"
	"strings"
	"time"
	"github.com/ktpswjz/httpserver/example/webserver/server/config"
	"github.com/ktpswjz/httpserver/example/webserver/server/controller"
)

type Doc struct {
	controller.Base
	Document document.Document
}

func (s *Doc) GetCatalogTree(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant)  {
	filter := &model.CatalogFilter {
		Keywords: "",
	}
	a.GetArgument(r, filter)

	a.Success(s.Document.GetCatalogTree(filter.Keywords))
}

func (s *Doc) GetFunction(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant)  {
	id := p.ByName("id")
	fun := s.Document.GetFunction(id)
	if fun == nil {
		a.Error(errors.NotExist)
		return
	}
	fun.FullPath = fmt.Sprintf("%s://%s%s", a.Schema(), r.Host, fun.Path)

	a.Success(fun)
}

func (s *Doc) CreateToken(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant)  {
	filter := &model.TokenFilter{
		Account: "",
	}
	err := a.GetArgument(r, filter)
	if err != nil {
		a.Error(errors.InputError,  err)
		return
	}
	account := strings.ToLower(strings.TrimSpace(filter.Account))
	if account == "" {
		a.Error(errors.InputInvalid,  "账号为空")
		return
	}

	var user *config.SiteAdminUser = nil
	userCount := len(s.Config.Site.Admin.Users)
	for index := 0; index < userCount; index++  {
		if account == strings.ToLower(s.Config.Site.Admin.Users[index].Account) {
			user = &s.Config.Site.Admin.Users[index]
			break
		}
	}

	if user == nil {
		a.Error(errors.LoginAccountNotExit,  err)
		return
	}

	now := time.Now()
	token := &model.Token{
		ID: a.GenerateGuid(),
		UserAccount: user.Account,
		LoginIP: a.RIP(),
		LoginTime: now,
		ActiveTime: now,
	}
	err = s.DbToken.Set(token)
	if err != nil {
		a.Error(errors.Exception,  err)
		return
	}

	login := &model.Login{
		Token: token.ID,
	}

	a.Success(login)
}