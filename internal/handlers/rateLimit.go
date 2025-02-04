package handlers

import (
	"context"
	"net/http"

	"github.com/mennanov/limiters"
	"github.com/the-redx/link-shortener/internal/services"
	"github.com/the-redx/link-shortener/pkg/errs"
	"go.uber.org/zap"
)

func RateLimitMW(next http.HandlerFunc, limiter services.RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value("Logger").(*zap.SugaredLogger)
		duration, err := limiter.Limit(context.TODO())

		if err == limiters.ErrLimitExhausted {
			logger.Debugf("Rate limit exceeded. Try again in %d seconds", int32(duration.Seconds()))
			writeError(w, errs.NewBadRequestError("Rate limit exceeded. Try later"))
			return
		} else if err != nil {
			logger.Debugf("Rate limiter error", err.Error())
			writeError(w, errs.NewUnexpectedError("Something went wrong"))
			return
		}

		next(w, r)
	}
}
