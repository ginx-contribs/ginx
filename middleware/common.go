package middleware

import (
	"github.com/246859/ginx/constant/status"
	"github.com/246859/ginx/contribs/accesslog"
	"github.com/246859/ginx/contribs/cors"
	"github.com/246859/ginx/contribs/ratelimit"
	"github.com/246859/ginx/contribs/recovery"
	"github.com/246859/ginx/pkg/resp"
	"github.com/gin-gonic/gin"
	"log/slog"
)

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

// NoRoute deals with case of 404
func NoRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp.New(ctx).Status(status.NotFound).Render()
	}
}

// NoMethod deals with case of method not allowed
func NoMethod() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp.New(ctx).Status(status.MethodNotAllowed).Render()
	}
}
