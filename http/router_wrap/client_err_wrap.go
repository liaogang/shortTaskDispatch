package router_wrap

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"root/http/writer"
	"root/manager"
	"root/page_wrap"
)

type ClientErrorRouterGroup struct {
	raw chi.Router
}

func NewClientErrorGroup(tmp chi.Router) *ClientErrorRouterGroup {
	var group = &ClientErrorRouterGroup{raw: tmp}

	return group
}

func (slf *ClientErrorRouterGroup) Post(pattern string, handler ClientErrorHandlerFunc) {
	slf.raw.Post(pattern, func(w http.ResponseWriter, r *http.Request) {

		pageId := r.URL.Query().Get("pageId")

		cli, err := manager.GetPageById(pageId)
		if err != nil {
			writer.Err0(w, err)
			return
		}

		clientErrorBridge(cli, w, r, handler)
	})
}

func (slf *ClientErrorRouterGroup) Get(pattern string, handler ClientErrorHandlerFunc) {
	slf.raw.Get(pattern, func(w http.ResponseWriter, r *http.Request) {

		pageId := r.URL.Query().Get("pageId")

		cli, err := manager.GetPageById(pageId)
		if err != nil {
			writer.Err0(w, err)
			return
		}

		clientErrorBridge(cli, w, r, handler)
	})
}

type ClientErrorHandlerFunc func(cli *page_wrap.Wrapper, w http.ResponseWriter, r *http.Request) error

func clientErrorBridge(cli *page_wrap.Wrapper, w http.ResponseWriter, r *http.Request, handler ClientErrorHandlerFunc) {
	var err = handler(cli, w, r)
	if err != nil {
		writer.Err0(w, err)
	}
}
