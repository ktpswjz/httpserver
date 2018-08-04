package handler

import (
	"net/http"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/router"
	"github.com/ktpswjz/httpserver/id"
	"github.com/ktpswjz/httpserver/security/rsakey"
	"time"
	"net"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func NewHandler(handle Handle, log types.Log, restart func() error) (Handler, error)  {
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

	if handle != nil {
		handle.Map(instance.router)
	}

	return instance, nil
}

type innerHandler struct {
	types.Base
	router *router.Router
	handle Handle
	requestId id.Generator
	sqlEntityId id.Generator
	randKey *rsakey.Private
	restart func() error
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

	if s.handle != nil {
		a.leaveTime = time.Now()
		go func(a *Assistant) {
			s.handle.PostRouting(a)
		}(a)
	}
}

func (s *innerHandler) newAssistant(w http.ResponseWriter, r *http.Request) *Assistant {
	instance := &Assistant {response: w, schema: "http"}
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

	return instance
}