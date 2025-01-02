package ollama

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter interface {
	Wait() error
}

type rateLimiter struct {
	limiter *rate.Limiter
}

func newRateLimiter(rps int) *rateLimiter {
	return &rateLimiter{
		limiter: rate.NewLimiter(rate.Every(time.Second/time.Duration(rps)), rps),
	}
}

func (r *rateLimiter) Wait() error {
	return r.limiter.Wait(context.Background())
}
