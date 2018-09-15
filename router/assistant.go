package router

import (
	"github.com/ktpswjz/httpserver/security/rsakey"
	"github.com/ktpswjz/httpserver/types"
	"net/http"
	"time"
)

type Assistant interface {
	CanUpdate() bool
	CanRestart() bool
	Restart() error
	GetBody(r *http.Request) ([]byte, error)
	GetArgument(r *http.Request, v interface{}) error
	GetXml(r *http.Request, v interface{}) error
	Success(data interface{})
	Error(err types.Error, errDetails ...interface{})
	OutputJson(code int, data interface{}, errSummary string, errDetails ...interface{})

	IsError() bool
	Set(key string, val interface{})
	Get(key string) (interface{}, bool)
	Del(key string) bool
	SetRecord(v bool)
	GetRecord() bool
	SetInput(v []byte)
	GetInput() []byte
	GetOutput() []byte
	GetParam() []byte
	Method() string
	Schema() string
	Path() string
	RID() uint64
	RIP() string
	EnterTime() time.Time
	LeaveTime() time.Time
	Token() string
	JsonWebToken() string
	ClientKey() *rsakey.Public
	RandKey() *rsakey.Private
	GenerateGuid() string
}
