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
