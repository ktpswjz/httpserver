package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ktpswjz/httpserver/id"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/security/rsakey"
	"github.com/ktpswjz/httpserver/types"
	"net"
	"net/http"
	"strings"
	"time"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func NewHandler(handle Handle, log types.Log, redirectToHttps bool, restart func() error) (Handler, error) {
	privateKey := &rsakey.Private{}
	err := privateKey.Create(1024)
	if err != nil {
		return nil, err
	}

	instance := &innerHandler{router: router.New(), handle: handle}
	instance.SetLog(log)
	instance.requestId = id.NewTime()
	instance.sqlEntityId = id.NewTime()
	instance.randKey = privateKey
	instance.restart = restart
	instance.redirectToHttps = redirectToHttps

	if handle != nil {
		handle.Map(instance.router)
	}

	return instance, nil
}

type innerHandler struct {
	types.Base
	router          *router.Router
	handle          Handle
	requestId       id.Generator
	sqlEntityId     id.Generator
	randKey         *rsakey.Private
	restart         func() error
	redirectToHttps bool
}

func (s *innerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	a := s.newAssistant(w, r)
	s.LogDebug("new request: rid=", a.rid,
		", rip=", a.rip,
		", host=", r.Host,
		", schema=", a.schema,
		", method=", r.Method,
		", path=", a.path)

	if s.redirectToHttps {
		if a.schema == "http" {
			if r.Method == "GET" {
				redirectUrl := fmt.Sprintf("https://%s%s", r.Host, a.path)
				http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
				return
			}
		}
	}

	defer func(a *Assistant) {
		s.postRouting(a)
	}(a)

	// 异常处理
	defer func(a *Assistant) {
		if err := recover(); err != nil {
			s.LogError(a.schema,
				" request error(rid=", a.rid,
				", path=", a.path,
				"rip=", a.rip,
				"): ", err)

			a.OutputJson(-1, nil, "internal exception", err)
		}
	}(a)

	if s.handle != nil {
		if s.handle.PreRouting(w, r, a) {
			return
		}
	}

	a.path = r.URL.Path
	s.router.ServeHTTP2(w, r, a)

	//if s.handle != nil {
	//	a.leaveTime = time.Now()
	//	go func(a *Assistant) {
	//		s.handle.PostRouting(a)
	//	}(a)
	//}
}

func (s *innerHandler) postRouting(a *Assistant) {
	// 异常处理
	defer func() {
		if err := recover(); err != nil {
			s.LogError("postRouting", err)
		}
	}()

	if s.handle != nil {
		a.leaveTime = time.Now()
		go func(a *Assistant) {
			// 异常处理
			defer func() {
				if err := recover(); err != nil {
					s.LogError("postRouting", err)
				}
			}()

			s.handle.PostRouting(a)
		}(a)
	}
}

func (s *innerHandler) newAssistant(w http.ResponseWriter, r *http.Request) *Assistant {
	instance := &Assistant{response: w, schema: "http"}
	instance.method = r.Method
	if r.TLS != nil {
		instance.schema = "https"
	}
	instance.keys = make(map[string]interface{})
	instance.record = false
	instance.enterTime = time.Now()
	instance.path = r.URL.Path
	instance.rid = s.requestId.New()
	instance.rip, _, _ = net.SplitHostPort(r.RemoteAddr)
	instance.randKey = s.randKey
	instance.restart = s.restart
	instance.token = r.Header.Get("token")
	if instance.token == "" {
		if r.Method == "GET" {
			instance.token = r.FormValue("token")
		}
	}
	if len(r.URL.Query()) > 0 {
		params := make([]*types.Query, 0)
		for k, v := range r.URL.Query() {
			param := &types.Query{Key: k}
			if len(v) > 0 {
				param.Value = v[0]
				if strings.ToLower(k) == "jwt" {
					instance.jwt = v[0]
				}
			}
			params = append(params, param)
		}
		instance.param, _ = json.Marshal(params)
	}

	return instance
}
