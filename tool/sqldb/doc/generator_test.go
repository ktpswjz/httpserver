package doc

import (
	"github.com/ktpswjz/database/sqldb"
	"os"
	"strings"
	"path/filepath"
	"github.com/ktpswjz/database/sqldb/mysql"
	"fmt"
	"testing"
	"github.com/ktpswjz/database/sqldb/mssql"
)

func TestGenerator_CreateWord(t *testing.T) {
	goPath := os.Getenv("GOPATH")
	paths := strings.Split(goPath, ";")
	if len(paths) > 1 {
		goPath = paths[0]
	}
	docPath := filepath.Join(goPath, "tmp","db.docx")
	t.Log("doc file: ", docPath)
	docFile, err := os.Create(docPath)
	if err != nil {
		t.Fatal(err)
	}
	defer docFile.Close()

	db := MySql()
	//db = SqlServer()
	generator := &Generator{Database: db}

	err = generator.CreateWord(docFile, "title")
	if err != nil {
		t.Fatal(err)
	}
}

func MySql() sqldb.SqlDatabase {
	goPath := os.Getenv("GOPATH")
	paths := strings.Split(goPath, ";")
	if len(paths) > 1 {
		goPath = paths[0]
	}
	cfgPath := filepath.Join(goPath, "tmp", "cfg", "httpserver_tool_mysql_test.json")
	cfg := &mysql.Connection{
		Server:   "172.0.0.1",
		Port:     3306,
		Schema:   "mysql",
		Charset:  "utf8",
		Timeout:  10,
		User:     "root",
		Password: "",
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

	return mysql.NewDatabase(cfg)
}

func SqlServer() sqldb.SqlDatabase {
	goPath := os.Getenv("GOPATH")
	paths := strings.Split(goPath, ";")
	if len(paths) > 1 {
		goPath = paths[0]
	}
	cfgPath := filepath.Join(goPath, "tmp", "cfg", "httpserver_tool_mssql_test.json")
	cfg := &mssql.Connection{
		Server:   "127.0.0.1",
		Port:     1433,
		Schema:   "test",
		Instance: "MSSQLSERVER",
		User:     "sa",
		Password: "",
		Timeout:  10,
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

	return mssql.NewDatabase(cfg)
}
