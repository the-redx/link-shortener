package handlers

import (
	"context"
	"net/http"

	"github.com/rs/xid"
	"github.com/the-redx/link-shortener/pkg/utils"
	"go.uber.org/zap"
)

var traceIDKey = "TraceID"
var loggerKey = "Logger"

func LogMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := xid.New().String()

		ctx := context.WithValue(
			r.Context(),
			traceIDKey,
			traceID,
		)

		log := utils.Logger.With(zap.String(traceIDKey, traceID))
		ctx = context.WithValue(ctx, loggerKey, log)

		log.Debugf("Request: %s %s", r.Method, r.URL.Path)

		w.Header().Add("X-Trace-ID", traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
