package recovery

import (
	"errors"
	"github.com/246859/ginx/pkg/resp"
	"github.com/246859/ginx/pkg/s2b"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

type Options struct {
	Logger  *slog.Logger
	Handler func(ctx *gin.Context, logger *slog.Logger, err any)
}

type Option func(options *Options)

func WithLogger(logger *slog.Logger) Option {
	return func(options *Options) {
		options.Logger = logger
	}
}

func WithHandler(handler func(ctx *gin.Context, logger *slog.Logger, err any)) Option {
	return func(options *Options) {
		options.Handler = handler
	}
}

// Recovery resolve the panic in current context
func Recovery(opts ...Option) gin.HandlerFunc {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	if options.Logger == nil {
		options.Logger = slog.Default()
	}

	if options.Handler == nil {
		options.Handler = func(ctx *gin.Context, logger *slog.Logger, err any) {
			options.Logger.ErrorContext(ctx, "[Panic Recovered]", slog.Any("error", err), slog.String("stack", s2b.BytesToString(debug.Stack())))
			resp.Fail(ctx).Status(http.StatusInternalServerError).Render()
			ctx.Abort()
		}
	}

	return func(ctx *gin.Context) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				var brokenPipe bool

				var err any
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				if ne, ok := panicErr.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne.Err, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				err = panicErr
				if !brokenPipe {
					options.Handler(ctx, options.Logger, err)
				} else {
					return
				}
			}
		}()

		ctx.Next()
	}
}
