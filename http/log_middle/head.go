package log_middle

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

func ChiLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t := time.Now()

		next.ServeHTTP(w, r)

		var build strings.Builder
		build.WriteString(r.Method)
		build.WriteByte(' ')
		build.WriteString(r.URL.String())
		build.WriteByte(' ')
		build.WriteString(time.Since(t).String())

		log.Info().Msg(build.String())
	})
}
