package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx/constant/status"
	"github.com/ginx-contribs/ginx/contribs/accesslog"
	"github.com/ginx-contribs/ginx/contribs/cache"
	"github.com/ginx-contribs/ginx/contribs/cors"
	"github.com/ginx-contribs/ginx/contribs/ratelimit"
	"github.com/ginx-contribs/ginx/contribs/recovery"
	"github.com/ginx-contribs/ginx/pkg/resp"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"slices"
	"time"
)

// NoRoute deals with case of 404
func NoRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp.New(ctx).Status(status.NotFound).JSON()
	}
}

// NoMethod deals with case of method not allowed
func NoMethod(methods ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !slices.Contains(methods, ctx.Request.Method) {
			resp.New(ctx).Status(status.MethodNotAllowed).JSON()
		}
	}
}

// Logger returns the access log handler
func Logger(logger *slog.Logger, msg string) gin.HandlerFunc {
	return accesslog.AccessLog(accesslog.DefaultOptions(logger, msg))
}

// Recovery returns the recovery handler
func Recovery(logger *slog.Logger, handler func(ctx *gin.Context, logger *slog.Logger, err any)) gin.HandlerFunc {
	return recovery.Recovery(recovery.WithLogger(logger), recovery.WithHandler(handler))
}

// CORS returns the cors handler
func CORS(options cors.Options) gin.HandlerFunc {
	return cors.New(options)
}

// RateLimit returns limiter handler
func RateLimit(limiter ratelimit.Limiter, errorHandler func(ctx *gin.Context, err error)) gin.HandlerFunc {
	return ratelimit.RateLimit(ratelimit.WithLimiter(limiter), ratelimit.WithErrorHandler(errorHandler))
}

// CacheMemory returns the cache in memory handler
func CacheMemory(prefix string, ttl time.Duration) gin.HandlerFunc {
	return cache.Cache(cache.WithTTL(ttl), cache.WithPrefix(prefix))
}

// CacheRedis returns the cache in redis handler
func CacheRedis(prefix string, ttl time.Duration, client *redis.Client) gin.HandlerFunc {
	return cache.Cache(cache.WithTTL(ttl), cache.WithPrefix(prefix), cache.WithStore(cache.NewRedisStore(client)))
}
