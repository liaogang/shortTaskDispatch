package router_wrap

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"root/http/writer"
)

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

type ErrorRouterGroup struct {
	raw chi.Router
}

func NewErrorGroup(tmp chi.Router) *ErrorRouterGroup {
	var group = &ErrorRouterGroup{raw: tmp}

	return group
}

func (slf *ErrorRouterGroup) Post(pattern string, handler ErrorHandlerFunc) {
	slf.raw.Post(pattern, func(w http.ResponseWriter, r *http.Request) {
		errorBridge(w, r, handler)
	})
}

func (slf *ErrorRouterGroup) Get(pattern string, handler ErrorHandlerFunc) {
	slf.raw.Get(pattern, func(w http.ResponseWriter, r *http.Request) {
		errorBridge(w, r, handler)
	})
}

func errorBridge(w http.ResponseWriter, r *http.Request, handler ErrorHandlerFunc) {
	var err = handler(w, r)
	if err != nil {
		_ = writer.Err1(w, err)
	}
}
