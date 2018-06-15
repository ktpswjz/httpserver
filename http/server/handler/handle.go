package handler

import (
	"github.com/ktpswjz/httpserver/router"
	"net/http"
)

type Handle interface {
	Map(router *router.Router)
	PreRouting(w http.ResponseWriter, r *http.Request, a router.Assistant) bool
	PostRouting(a router.Assistant)
}