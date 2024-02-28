package log_middle

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()

		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		ctx.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		var event *zerolog.Event
		var xErr = ctx.Errors.Last()
		if xErr != nil {
			event = log.Warn()
		} else {
			event = log.Info()
		}

		if len(ctx.Errors) > 0 {
			event.Err(xErr)
		}

		var build strings.Builder
		build.WriteString(ctx.Request.Method)
		build.WriteByte(' ')
		build.WriteString(path)
		build.WriteString(" -> ")
		build.WriteString(http.StatusText(ctx.Writer.Status()))
		build.WriteByte(' ')
		build.WriteString(time.Since(t).String())

		//proxy, _ := http_util.ProxyFromReqHeader(ctx)
		//if proxy != nil {
		//	build.WriteByte(' ')
		//	build.WriteByte('[')
		//	build.WriteString(proxy.Address)
		//	build.WriteByte('@')
		//	build.WriteString(proxy.User)
		//	build.WriteByte(':')
		//	build.WriteString(proxy.Password)
		//	build.WriteByte(']')
		//}

		event.Msg(build.String())
	}
}
