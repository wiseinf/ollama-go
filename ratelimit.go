package ollama

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter 接口定义速率限制行为
type RateLimiter interface {
	Wait() error
}

// rateLimiter 实现令牌桶算法的速率限制
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
