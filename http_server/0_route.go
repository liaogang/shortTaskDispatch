package http_server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"root/http/log_middle"
	"root/http/router_wrap"
)

func StartRouter(address string) error {

	h := chi.NewRouter()

	h.Use(log_middle.ChiLogger)
	h1 := router_wrap.NewErrorGroup(h)

	h1.Post("/server_publish_task_and_wait_finish", publishTask)

	h1.Post("/client_claim_any_task", claimTask)

	h1.Post("/client_finish_task_with_success", finishTaskWithSuccess)
	h1.Post("/client_finish_task_with_error", finishTaskWithError)

	return http.ListenAndServe(address, h)
}
