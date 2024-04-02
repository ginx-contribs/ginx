package bbr

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"
)

type Limiter struct {
	limiter *bbr.BBR
}

func (l Limiter) Allow(ctx *gin.Context) (func(), error) {
	done, err := l.limiter.Allow()
	return func() { done(ratelimit.DoneInfo{}) }, err
}

func NewLimiter(bbr *bbr.BBR) *Limiter {
	return &Limiter{limiter: bbr}
}
