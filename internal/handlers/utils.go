package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/the-redx/link-shortener/pkg/errs"
)

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func writeError(w http.ResponseWriter, appErr *errs.AppError) {
	writeResponse(w, appErr.Code, appErr)
}

func getUserIdFromContext(r *http.Request) (string, *errs.AppError) {
	userId, ok := r.Context().Value("UserID").(string)
	if !ok {
		return "", errs.NewForbiddenError("Authentication error")
	}

	return userId, nil
}
