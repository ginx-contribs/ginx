package bucket

import (
	"github.com/gin-gonic/gin"
	ginxratelmit "github.com/ginx-contribs/ginx/contribs/ratelimit"
	"github.com/juju/ratelimit"
	"time"
)

// Limiter implements ratelimit.RateLimit interface by Token Bucket.
type Limiter struct {
	bucket *ratelimit.Bucket
	// max timeout to wait buckets becomes available
	maxWait time.Duration
}

func (b *Limiter) Allow(ctx *gin.Context) (func(), error) {
	var take int64
	if b.maxWait <= 0 {
		if b.bucket.TakeAvailable(take) <= 0 {
			return nil, ginxratelmit.ErrRateLimitExceed
		}
	} else {
		if !b.bucket.WaitMaxDuration(take, b.maxWait) {
			return nil, ginxratelmit.ErrRateLimitExceed
		}
	}
	return func() {}, nil
}

type Option func(limiter *Limiter)

func WithBucket(bucket *ratelimit.Bucket) Option {
	return func(limiter *Limiter) {
		limiter.bucket = bucket
	}
}

func WithMaxWait(wait time.Duration) Option {
	return func(limiter *Limiter) {
		limiter.maxWait = wait
	}
}

func NewLimiter(opts ...Option) *Limiter {
	var limiter Limiter
	for _, opt := range opts {
		opt(&limiter)
	}

	if limiter.bucket == nil {
		limiter.bucket = ratelimit.NewBucket(time.Second, 100)
	}

	return &limiter
}
