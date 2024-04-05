package ratelimit

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx/constant/status"
	"github.com/ginx-contribs/ginx/pkg/resp"
)

type Options struct {
	Limiter      Limiter
	ErrorHandler func(ctx *gin.Context, err error)
}

type Option func(options *Options)

func WithLimiter(limiter Limiter) Option {
	return func(options *Options) {
		options.Limiter = limiter
	}
}

func WithErrorHandler(errorHandler func(ctx *gin.Context, err error)) Option {
	return func(options *Options) {
		options.ErrorHandler = errorHandler
	}
}

// RateLimit returns a new limiter handler with options
func RateLimit(opts ...Option) gin.HandlerFunc {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	if options.Limiter == nil {
		options.Limiter = defaultLimiter{}
	}

	if options.ErrorHandler == nil {
		options.ErrorHandler = func(ctx *gin.Context, err error) {
			if err != nil {
				if errors.Is(err, ErrRateLimitExceed) { // rate limit exceeded
					resp.Fail(ctx).Status(status.TooManyRequests).Error(err).JSON()
				} else { // internal server error
					resp.Fail(ctx).Status(status.InternalServerError).Error(err).JSON()
				}
				ctx.Abort()
			}
		}
	}

	return func(ctx *gin.Context) {
		// try to allow request
		done, err := options.Limiter.Allow(ctx)
		// allow
		if err == nil {
			ctx.Next()
			done()
		} else if err != nil {
			options.ErrorHandler(ctx, err)
		}
	}
}
