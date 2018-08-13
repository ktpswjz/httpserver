package handler

import (
	"net/http"
	"github.com/ktpswjz/httpserver/security/rsakey"
	"time"
	"github.com/ktpswjz/httpserver/types"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/ktpswjz/httpserver/id"
	"runtime"
)

type Assistant struct {
	response http.ResponseWriter
	method string
	schema string
	path string
	rid uint64
	rip string
	token string
	jwt string
	clientKey *rsakey.Public
	randKey *rsakey.Private
	restart func() error

	keys map[string]interface{}
	record bool
	input []byte
	output []byte
	param []byte
	outputCode *int
	enterTime time.Time
	transferTime time.Time
	leaveTime time.Time
}

func (s *Assistant) CanUpdate() bool {
	if s.restart == nil {
		return false
	} else if runtime.GOOS == "linux" {
		return true
	} else {
		return false
	}
}

func (s *Assistant) CanRestart() bool {
	if s.restart == nil {
		return false
	} else {
		return true
	}
}

func (s *Assistant) Restart() error {
	if s.restart == nil {
		return fmt.Errorf("restart not supported")
	}

	return s.restart()
}

func (s *Assistant) Success(data interface{}) {
	s.OutputJson(0, data, "")
}

func (s *Assistant) Error(err types.Error, errDetails ...interface{}) {
	s.OutputJson(err.Code(), nil, err.Summary(), errDetails...)
}

func (s *Assistant) OutputJson(code int, data interface{}, errSummary string, errDetails ...interface{}) {
	if s.response == nil {
		return
	}
	s.outputCode = &code

	result := &types.Result{
		Code: code,
		Data: data,
		Elapse: time.Now().Sub(s.enterTime).String(),
		Serial: s.rid,
		Error: types.ResultError {
			Summary: errSummary,
			Detail: fmt.Sprint(errDetails...),
		},
	}

	resultData, err := result.Marshal()
	s.transferTime = time.Now()
	if err != nil {
		fmt.Fprint(s.response, err)
	} else {
		s.response.Header().Add("Access-Control-Allow-Origin", "*")
		s.response.Header().Set("Content-Type", "application/json;charset=utf-8")
		s.response.Write(resultData)
		s.output = resultData
	}
}

func (s *Assistant) GetBody(r *http.Request) ([]byte, error) {
	return ioutil.ReadAll(r.Body)
}

func (s *Assistant) GetArgument(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err == nil {
		s.input, _ = json.Marshal(v)
	}

	return err
}

func (s *Assistant) IsError() bool  {
	if s.outputCode == nil {
		return false
	}

	if *s.outputCode == 0 {
		return false
	}

	return true
}

func (s *Assistant) Set(key string, val interface{})  {
	s.keys[key] = val
}

func (s *Assistant) Get(key string) (interface{}, bool)  {
	val, ok := s.keys[key]
	if ok {
		return val, true
	} else {
		return nil, false
	}
}

func (s *Assistant) Del(key string) bool  {
	_, ok := s.keys[key]
	if ok {
		delete(s.keys, key)
		return true
	} else {
		return false
	}
}

func (s *Assistant) SetRecord(v bool) {
	s.record = v
}

func (s *Assistant) GetRecord() bool  {
	return s.record
}

func (s *Assistant) SetInput(v []byte)  {
	s.input = v
}

func (s *Assistant) GetInput() []byte  {
	return s.input
}

func (s *Assistant) GetOutput() []byte  {
	return s.output
}

func (s *Assistant) GetParam() []byte  {
	return s.param
}

func (s *Assistant) Method() string {
	return s.method
}

func (s *Assistant) Schema() string {
	return s.schema
}

func (s *Assistant) Path() string {
	return s.path
}
func (s *Assistant) RID() uint64  {
	return s.rid
}

func (s *Assistant) RIP() string  {
	return s.rip
}
func (s *Assistant) EnterTime() time.Time  {
	return s.enterTime
}
func (s *Assistant) LeaveTime() time.Time  {
	return s.leaveTime
}

func (s *Assistant) Token() string  {
	return s.token
}

func (s *Assistant) JsonWebToken() string  {
	return s.jwt
}

func (s *Assistant) ClientKey() *rsakey.Public  {
	return s.clientKey
}

func (s *Assistant) RandKey() *rsakey.Private {
	return s.randKey
}

func (s *Assistant) GenerateGuid() string {
	return id.GenerateGuid()
}