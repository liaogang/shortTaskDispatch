package http_server

import (
	"github.com/gin-gonic/gin"
	"root/http_server/internal/log_middle"
)

func SetupRouter(address string) {
	gin.SetMode(gin.ReleaseMode)
	h := gin.New()
	h.Use(log_middle.Logger())

	//发布任务请求，并等待完成
	h.POST("publishTask", genHandler(publishTask))

	//认领任务
	h.POST("claimTask", genHandler(claimTask))

	//完成认领的任务
	h.POST("finishTask", genHandler(finishTask))

	h.Run(address)
}
