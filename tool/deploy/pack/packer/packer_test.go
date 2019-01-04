package packer

import (
	"fmt"
	"github.com/ktpswjz/httpserver/tool/deploy/pack/config"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPacker_Pack(t *testing.T) {
	pack := NewPacker(packConfig())
	err := pack.Pack()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("成功,输出目录:", pack.OutputFolder())
}

func packConfig() *config.Config {
	goPath := os.Getenv("GOPATH")
	paths := strings.Split(goPath, ";")
	if len(paths) > 1 {
		goPath = paths[0]
	}
	cfgPath := filepath.Join(goPath, "tmp", "cfg", "httpserver_tool_deploy_pack_test.json")
	cfg := &config.Config{
		Version:     "1.0.1.0",
		Source:      true,
		Destination: filepath.Join(goPath, "tmp", "rel"),
		App: config.App{
			Name: "epes",
			Bin: config.Binary{
				Root: filepath.Join(goPath, "tmp", "bin"),
				Files: map[string]string{
					"go_build_epes_service": "epes",
				},
			},
			Src: config.Source{
				Root: filepath.Join(goPath, "src", "epes"),
				Ignore: []string{
					"tool",
					".git",
					".idea",
					".gitignore",
					"README.md",
				},
			},
		},
		Webs: []config.Web{
			{
				Enable: false,
				Name:   "staff",
				Src: config.Source{
					Root: "/home/dev/vue/epes/staff",
					Ignore: []string{
						"node_modules",
						"dist",
						".git",
						".idea",
						".gitignore",
						"README.md",
					},
				},
			},
			{
				Enable: true,
				Name:   "wcp",
				Src: config.Source{
					Root: "/home/dev/vue/epes/wechat/public",
					Ignore: []string{
						"node_modules",
						"dist",
						".git",
						".idea",
						".gitignore",
						"README.md",
					},
				},
			},
			{
				Enable: true,
				Name:   "operation",
				Src: config.Source{
					Root: "/home/dev/vue/epes/operation",
					Ignore: []string{
						"node_modules",
						"dist",
						".git",
						".idea",
						".gitignore",
						"README.md",
					},
				},
			},
			{
				Enable: true,
				Name:   "document",
				Src: config.Source{
					Root: "/home/dev/vue/epes/document",
					Ignore: []string{
						"node_modules",
						"dist",
						".git",
						".idea",
						".gitignore",
						"README.md",
					},
				},
			},
		},
	}
	_, err := os.Stat(cfgPath)
	if os.IsNotExist(err) {
		err = cfg.SaveToFile(cfgPath)
		if err != nil {
			fmt.Println("generate configure file fail: ", err)
		}
	} else {
		err = cfg.LoadFromFile(cfgPath)
		if err != nil {
			fmt.Println("load configure file fail: ", err)
		}
	}

	return cfg
}
