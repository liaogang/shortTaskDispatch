package http_server

import "github.com/gin-gonic/gin"

type handlerReturnErr func(ctx *gin.Context) error

func genHandler(h1 handlerReturnErr) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err = h1(ctx)

		if err != nil {
			ginWriteErr(ctx, err)
		}

	}
}

func ginWriteErr(ctx *gin.Context, err error) {

	ctx.Error(err)

	ctx.Header("X-Error", "1")
	ctx.Header("Content-Type", "plain/txt; charset=utf-8")
	ctx.Writer.Write([]byte(err.Error()))
}
