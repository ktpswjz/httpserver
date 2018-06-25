package admin

import (
	"github.com/ktpswjz/httpserver/example/webserver/server/controller"
	"net/http"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"bytes"
	"path/filepath"
	"os"
	"fmt"
	"github.com/ktpswjz/httpserver/archive"
	"github.com/ktpswjz/httpserver/example/webserver/model"
	"github.com/ktpswjz/httpserver/types"
	"time"
	"github.com/ktpswjz/httpserver/document"
	"github.com/ktpswjz/httpserver/example/webserver/server/config"
	"github.com/ktpswjz/httpserver/example/webserver/database/memory"
)

type App struct {
	controller.Base
}

func NewApp(cfg *config.Config, log types.Log, dbToken memory.Token) *App  {
	instance := &App{}
	instance.Config = cfg
	instance.SetLog(log)
	instance.DbToken = dbToken

	return instance
}

func (s *App) Upload(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	appPath := r.FormValue("path")
	if appPath == "" {
		a.Error(errors.InputInvalid, "path name is empty")
		return
	}

	appFile, _, err := r.FormFile("file")
	if err != nil {
		a.Error(errors.InputInvalid, "invalid file: ", err)
		return
	}
	defer appFile.Close()
	var buf bytes.Buffer
	fileSize, err := buf.ReadFrom(appFile)
	if err != nil {
		a.Error(errors.InputInvalid, "read file error: ", err)
		return
	}
	if fileSize < 0 {
		a.Error(errors.InputInvalid, "invalid file: size is zero")
		return
	}

	tempFolder := filepath.Join(s.Config.Site.App.Root, a.GenerateGuid())
	err = os.MkdirAll(tempFolder, 0777)
	if err != nil {
		a.Error(errors.InputInvalid, fmt.Sprintf("create temp folder '%s' error:", tempFolder), err)
		return
	}
	defer os.RemoveAll(tempFolder)

	fileData := buf.Bytes()
	zipFile := &archive.Zip{}
	err = zipFile.DecompressMemory(fileData, tempFolder)
	if err != nil {
		a.Error(errors.InputInvalid, "decompress file error: ", err)
		return
	}

	appFolder := filepath.Join(s.Config.Site.App.Root, appPath)
	err = os.RemoveAll(appFolder)
	if err != nil {
		a.Error(errors.InputInvalid, fmt.Sprintf("remove original app '%s' error:", appPath), err)
		return
	}
	os.MkdirAll(filepath.Dir(appFolder), 0777)
	err = os.Rename(tempFolder, appFolder)
	if err != nil {
		a.Error(errors.InputInvalid, fmt.Sprintf("rename app folder '%s' error:", appFolder), err)
		return
	}
	appInfo := &model.App{
		UploadTime: types.Time(time.Now()),
	}
	appInfo.Remark = r.FormValue("remark")
	appInfo.Version = r.FormValue("version")
	if appInfo.Version == "" {
		appInfo.Version = "1.0.1.0"
	}
	token, err := s.DbToken.Get(a.Token())
	if err == nil && token != nil {
		appInfo.UploadUser = token.UserAccount
	}
	appInfoName  := filepath.Join(appFolder, "app.info")
	appInfo.SaveToFile(appInfoName)

	a.Success(appPath)
}

func (s *App) Tree(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	baseUrl := fmt.Sprintf("%s://%s/app", a.Schema(), r.Host)
	appTree := &model.AppTree{}
	appTree.ParseChildren(s.Config.Site.App.Root, baseUrl)

	a.Success(appTree.Children)
}

func (s *App) TreeDoc(a document.Assistant) document.Function  {
	function := a.CreateFunction("获取应用程序列表")
	function.SetNote("获取服务系统当前所有应用程序")
	function.SetOutputExample( []model.AppTree {
		{
			Path:       "test",
			UploadTime: types.Time(time.Now()),
			UploadUser: "admin",
			Version:    "1.0.1.0",
			Url:        "https://www.example.com/app/test",
		},
	})
	function.SetContentType("")

	catalog := a.CreateCatalog("平台管理", "平台管理服务相关接口")
	catalog.CreateChild("应用程序", "应用程序相关接口").SetFunction(function)

	return function
}

func (s *App) Delete(w http.ResponseWriter, r *http.Request, p router.Params, a router.Assistant) {
	filter := &model.AppFilter{}
	err := a.GetArgument(r, filter)
	if err != nil {
		a.Error(errors.InputError,  err)
		return
	}
	if filter.Path == "" {
		a.Error(errors.InputError,  "应用程序路径为空")
		return
	}

	appPath := filepath.Join(s.Config.Site.App.Root, filter.Path)
	info, err := os.Stat(appPath)
	if os.IsNotExist(err) {
		a.Error(errors.InputError,  fmt.Sprintf("应用程序'%s'不存在", filter.Path))
		return
	}
	if !info.IsDir() {
		a.Error(errors.InputError,  fmt.Sprintf("应用程序'%s'无效", filter.Path))
		return
	}
	err = os.RemoveAll(appPath)
	if err != nil {
		a.Error(errors.Exception,  err)
		return
	}

	a.Success(filter.Path)
}

func (s *App) DeleteDoc(a document.Assistant) document.Function  {
	function := a.CreateFunction("删除应用程序")
	function.SetNote("删除存在的应用程序")
	function.SetInputExample(&model.AppFilter{
		Path: "test",
	})
	function.SetOutputExample("test")

	catalog := a.CreateCatalog("平台管理", "平台管理服务相关接口")
	catalog.CreateChild("应用程序", "应用程序相关接口").SetFunction(function)

	return function
}
