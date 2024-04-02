package ratelimit

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var ErrRateLimitExceed = errors.New("rate limit exceeded")

// Limiter limits rate of requests for server
type Limiter interface {
	Allow(ctx *gin.Context) (func(), error)
}

// defaultLimiter do nothing
type defaultLimiter struct{}

func (d defaultLimiter) Allow(ctx *gin.Context) (func(), error) {
	return func() {}, nil
}
