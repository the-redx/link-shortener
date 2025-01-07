package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/the-redx/link-shortener/pkg/errs"
	"go.uber.org/zap"
)

var userIdKey = "UserID"

func AuthMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value("Logger").(*zap.SugaredLogger)

		userId := r.Header.Get("X-User-ID")
		userId = strings.TrimSpace(userId)

		logger.Debug("Auth middleware", zap.String("userId", userId))

		if userId == "" {
			logger.Debug("Authentication error")
			writeError(w, errs.NewForbiddenError("Authentication error"))
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, userId)
		ctx = context.WithValue(ctx, "Logger", logger.With(zap.String("UserID", userId)))

		next(w, r.WithContext(ctx))
	}
}
