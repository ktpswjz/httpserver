package config

import (
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/http/server/configure"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	version *types.Version

	Name 	string 					`json:"name"`
	Log 	Log 					`json:"log"`
	Server 	configure.Server		`json:"server"`
	Site 	Site 					`json:"site"`
}


func NewConfig() *Config  {
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
					File: "",
					Password: "",
				},
			},

		},
		Site: Site{
			App: SiteApp{
				Root: "",
			},
			Doc: SiteDoc{
				Enable: true,
				Root: "/home/dev/project/vue/httpserver/doc/dist",
			},
			Admin: SiteAdmin{
				Enable: true,
				Root: "/home/dev/project/vue/httpserver/admin/dist",
				Api: SiteAdminApi{
					Token: Token{
						Expiration: 30,
					},
				},
				Users: []SiteAdminUser {
					{
						Account: "admin",
						Password: "1",
					},
				},
				Ldap: SiteAdminLdap {
					Enable: true,
					Host: "192.168.123.1",
					Port: 389,
					Base: "dc=csby,dc=studio",
				},
			},
		},
	}
}

func (s *Config) SetVersion(version *types.Version)  {
	s.version = version
}

func (s *Config) GetVersion() *types.Version  {
	return s.version
}

func (s *Config) GetServer() *configure.Server  {
	return &s.Server
}

func (s *Config) LoadFromFile(filePath string) error {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}

func (s *Config) SaveToFile(filePath string) error {
	return nil
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
