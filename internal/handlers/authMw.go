package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/the-redx/link-shortener/pkg/errs"
)

func AuthMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("X-User-ID")
		userId = strings.TrimSpace(userId)

		if userId == "" {
			writeError(w, errs.NewForbiddenError("Authentication error"))
			return
		}

		ctx := context.WithValue(r.Context(), "UserID", userId)
		next(w, r.WithContext(ctx))
	}
}
