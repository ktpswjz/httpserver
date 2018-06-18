package router

import (
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/security/rsakey"
	"net/http"
)

type Assistant interface {
	Restart() error
	GetBody(r *http.Request) ([]byte, error)
	GetArgument(r *http.Request, v interface{}) error
	Success(data interface{})
	Error(err types.Error, errDetails ...interface{})
	OutputJson(code int, data interface{}, errSummary string, errDetails ...interface{})

	SetRecord(v bool)
	GetRecord() bool
	SetInput(v []byte)
	GetInput() []byte
	GetOutput() []byte
	Schema() string
	RID() uint64
	RIP() string
	Token() string
	ClientKey() *rsakey.Public
	RandKey() *rsakey.Private
	GenerateGuid() string
}