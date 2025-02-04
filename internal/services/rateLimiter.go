package services

import (
	"context"
	"time"

	"github.com/mennanov/limiters"
)

type RateLimiter interface {
	Limit(context context.Context) (time.Duration, error)
}

func NewRateLimiter(capacity int64, rate time.Duration) RateLimiter {
	limiter := limiters.NewSlidingWindow(
		capacity,
		rate,
		limiters.NewSlidingWindowInMemory(),
		limiters.NewSystemClock(),
		0.01,
	)

	return limiter
}
