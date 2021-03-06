package config

import (
	"encoding/json"
	"fmt"
	"github.com/ktpswjz/httpserver/http/server/configure"
	"github.com/ktpswjz/httpserver/types"
	"io/ioutil"
	"os"
	"sync"
)

type Config struct {
	mutex sync.RWMutex
	args  *types.Args

	Name   string           `json:"name"`
	Log    Log              `json:"log"`
	Server configure.Server `json:"server"`
	Site   Site             `json:"site"`
}

func NewConfig() *Config {
	return &Config{
		Name: "HTTP服务器",
		Log: Log{
			Folder: "",
			Level:  "error|warning|info|debug",
		},
		Server: configure.Server{
			Http: configure.Http{
				Enabled: true,
				Address: "",
				Port:    "8080",
			},
			Https: configure.Https{
				Enabled: true,
				Address: "",
				Port:    "8443",
				Cert: configure.Certificate{
					File:     "",
					Password: "",
				},
			},
		},
		Site: Site{
			App: SiteApp{
				Root: "/home/dev/project/apps",
			},
			Doc: SiteDoc{
				Enable: true,
				Root:   "/home/dev/project/vue/httpserver/doc/dist",
			},
			Admin: SiteAdmin{
				Enable: true,
				Root:   "/home/dev/project/vue/httpserver/admin/dist",
				Api: SiteAdminApi{
					Token: Token{
						Expiration: 30,
					},
				},
				Users: []SiteAdminUser{
					{
						Account:  "admin",
						Password: "1",
					},
				},
				Ldap: SiteAdminLdap{
					Enable: true,
					Host:   "192.168.123.1",
					Port:   389,
					Base:   "dc=csby,dc=studio",
				},
			},
		},
	}
}

func (s *Config) SetArgs(args *types.Args) {
	s.args = args
}

func (s *Config) GetArgs() *types.Args {
	return s.args
}

func (s *Config) GetServer() *configure.Server {
	return &s.Server
}

func (s *Config) LoadFromFile(filePath string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}

func (s *Config) SaveToFile(filePath string) error {
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

func (s *Config) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(bytes[:])
}

func (s *Config) FormatString() string {
	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return ""
	}

	return string(bytes[:])
}
