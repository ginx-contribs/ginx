package counter

import (
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx/contribs/ratelimit"
	"golang.org/x/net/context"
	"time"
)

// Counter counts frequency of specific key whether is go beyond the max limit in window period
type Counter interface {
	Count(ctx context.Context, key string, limit int, window time.Duration) (int, error)
}

// Limiter implements ratelimit.Limiter interface by Counter
type Limiter struct {
	Limit   int
	Window  time.Duration
	KeyFn   func(ctx *gin.Context) string
	Counter Counter
}

func (c *Limiter) Allow(ctx *gin.Context) (func(), error) {
	key := c.KeyFn(ctx)
	count, err := c.Counter.Count(ctx, key, c.Limit, c.Window)
	if err != nil {
		return nil, err
	}
	if count >= c.Limit {
		return nil, ratelimit.ErrRateLimitExceed
	}
	return func() {}, nil
}

type Option func(options *Limiter)

func WithLimit(limit int) Option {
	return func(cl *Limiter) {
		cl.Limit = limit
	}
}

func WithWindow(window time.Duration) Option {
	return func(cl *Limiter) {
		cl.Window = window
	}
}

func WithKeyFn(keyFn func(ctx *gin.Context) string) Option {
	return func(options *Limiter) {
		options.KeyFn = keyFn
	}
}

func WithCounter(counter Counter) Option {
	return func(cl *Limiter) {
		cl.Counter = counter
	}
}

func ClientIpKey() func(ctx *gin.Context) string {
	return func(ctx *gin.Context) string {
		return ctx.ClientIP()
	}
}

func UrlKey() func(ctx *gin.Context) string {
	return func(ctx *gin.Context) string {
		return ctx.Request.RequestURI
	}
}

func NewLimiter(options ...Option) *Limiter {
	var limiter Limiter
	for _, option := range options {
		option(&limiter)
	}

	if limiter.Counter == nil {
		limiter.Counter = Cache()
	}

	if limiter.KeyFn == nil {
		limiter.KeyFn = UrlKey()
	}

	if limiter.Limit == 0 {
		limiter.Limit = 100
	}

	if limiter.Window == 0 {
		limiter.Window = time.Minute
	}

	return &limiter
}
