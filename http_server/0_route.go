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

//func resolveCaptcha(ctx *gin.Context) error {
//
//	url := ctx.PostForm("url")
//
//	ctx.PostForm("cookie")
//
//	ctx.PostForm("proxy_address")
//	ctx.PostForm("proxy_user")
//	ctx.PostForm("proxy_password")
//
//	ctx.PostForm("timeout")
//
//	//发放任务
//	_ = url
//
//	//收回任务
//
//	return nil
//}
