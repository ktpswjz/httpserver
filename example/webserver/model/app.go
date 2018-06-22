package model

import (
	"github.com/ktpswjz/httpserver/types"
	"sync"
	"io/ioutil"
	"encoding/json"
	"os"
	"fmt"
	"path/filepath"
)

type App struct {
	mutex 	sync.RWMutex

	Version		string		`json:"version" note:"版本号"`
	UploadTime 	types.Time	`json:"uploadTime" note:"上传时间"`
	UploadUser	string		`json:"uploadUser" note:"上传者账号"`
	Remark		string		`json:"remark" note:"说明"`
}

func (s *App) LoadFromFile(filePath string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}

func (s *App) SaveToFile(filePath string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, string(bytes[:]))

	return err
}

type AppFilter struct {
	Name		string		`json:"name" note:"应用程序名称"`
}

type AppTree struct {
	parent		*AppTree

	Type		int			`json:"type" note:"0-folder; 1-app"`
	Name		string		`json:"name" note:"应用程序名称"`
	Path		string		`json:"path" note:"应用程序路径"`
	Url			string		`json:"url" note:"访问地址"`
	Version		string		`json:"version" note:"版本号"`
	UploadTime 	types.Time	`json:"uploadTime" note:"上传时间"`
	UploadUser	string		`json:"uploadUser" note:"上传者账号"`
	Remark		string		`json:"remark" note:"说明"`

	Children	[]*AppTree	`json:"children"`
}

func (s *AppTree) ParseChildren(folderPath, baseUrl string)  {
	s.Children = make([]*AppTree, 0)

	paths, err := ioutil.ReadDir(folderPath)
	if err == nil {
		for _, path := range paths {
			if !path.IsDir() {
				continue
			}

			child := &AppTree {
				Name: path.Name(),
				parent: s,
				Type: 0,
				UploadTime: types.Time(path.ModTime()),
			}
			child.Path = child.getPath()
			s.Children = append(s.Children, child)

			appInfo  := filepath.Join(folderPath, path.Name(), "app.info")
			app := &App {
				Version: "1.0.1.0",
				UploadTime: types.Time(path.ModTime()),
			}
			err = app.LoadFromFile(appInfo)
			if err == nil {
				child.Type = 1
				child.Version = app.Version
				child.UploadUser = app.UploadUser
				child.Remark = app.Remark
				child.Url = fmt.Sprintf("%s/%s", baseUrl, child.Path)
				child.Children = make([]*AppTree, 0)
			} else {
				child.ParseChildren(filepath.Join(folderPath, path.Name()), baseUrl)
			}
		}
	}
}

func (s *AppTree) getPath() string  {
	path := s.Name

	parent := s.parent
	for parent != nil {
		if parent.Name == "" {
			break
		}

		path = parent.Name + "/" + path
		parent = parent.parent
	}

	return path
}