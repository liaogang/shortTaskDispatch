package http_server

import (
	"github.com/go-chi/chi/v5"
	"root/http/log_middle"
	"root/http/router_wrap"
)

func StartRouter(address string) error {

	h := chi.NewRouter()

	h.Use(log_middle.ChiLogger)
	h1 := router_wrap.NewErrorGroup(h)

	h1.Post("/server_publish__task_and_wait_finish", publishTask)

	h1.Post("/client_claim_any_task", claimTask)

	h1.Post("/client_finish_task_with_success", finishTaskWithSuccess)
	h1.Post("/client_finish_task_with_error", finishTaskWithError)

	return nil
}

//func SetupRouter(address string) error {
//	gin.SetMode(gin.ReleaseMode)
//	h := gin.New()
//	h.Use(log_middle.Logger())
//
//	//发布任务请求，并等待完成
//	h.POST("publishTask", genHandler(publishTask))
//
//	//认领任务
//	h.POST("claimTask", genHandler(claimTask))
//
//	//完成认领的任务
//	h.POST("finishTask", genHandler(finishTask))
//
//	h.Run(address)
//
//	return nil
//}
