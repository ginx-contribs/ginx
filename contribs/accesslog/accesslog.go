package accesslog

import (
	"fmt"
	"github.com/246859/ginx/constant/headers"
	"github.com/dstgo/size"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http/httputil"
	"time"
)

func DefaultOptions(logger *slog.Logger, msg string) Options {
	return Options{
		Logger:           logger,
		Msg:              msg,
		ShowCost:         true,
		ShowIp:           true,
		ShowAgent:        true,
		ShowURL:          true,
		ShowPath:         true,
		ShowRoute:        true,
		ShowRequestId:    true,
		ShowRequestSize:  true,
		ShowResponseSize: true,
		ShowError:        true,
	}
}

type Options struct {
	Logger           *slog.Logger
	Msg              string
	ShowCost         bool
	ShowIp           bool
	ShowAgent        bool
	ShowURL          bool
	ShowPath         bool
	ShowRoute        bool
	ShowRequestId    bool
	ShowRequestSize  bool
	ShowResponseSize bool
	ShowError        bool
}

// AccessLog records server access logs
func AccessLog(options Options) gin.HandlerFunc {
	if options.Msg == "" {
		options.Msg = "access logs"
	}
	logger := options.Logger

	return func(ctx *gin.Context) {
		begin := time.Now()
		ctx.Next()

		var (
			status = ctx.Writer.Status()
			method = ctx.Request.Method
			cost   = time.Now().Sub(begin)
		)

		attrs := []any{slog.Int("status", status), slog.String("method", method)}

		if options.ShowCost {
			attrs = append(attrs, slog.String("cost", roundDuration(cost)))
		}

		if options.ShowIp {
			attrs = append(attrs, slog.String("ip", ctx.ClientIP()))
		}

		if options.ShowAgent {
			attrs = append(attrs, slog.String("agent", ctx.Request.UserAgent()))
		}

		if options.ShowURL {
			attrs = append(attrs, slog.String("url", ctx.Request.RequestURI))
		}

		if options.ShowPath {
			attrs = append(attrs, slog.String("path", ctx.Request.URL.Path))
		}

		if options.ShowRoute {
			attrs = append(attrs, slog.String("route", ctx.FullPath()))
		}

		if options.ShowRequestSize {
			req, _ := httputil.DumpRequest(ctx.Request, true)
			attrs = append(attrs, slog.String("request-size", roundSize(len(req)).String()))
		}

		if options.ShowResponseSize {
			attrs = append(attrs, slog.String("response-size", roundSize(int(ctx.Writer.Size())).String()))
		}

		if options.ShowRequestId {
			requestId := ctx.Writer.Header().Get(headers.XRequestId)
			if requestId != "" {
				attrs = append(attrs, slog.String("request-id", requestId))
			}
		}

		if options.ShowError && len(ctx.Errors) != 0 {
			attrs = append(attrs, slog.String("errors", ctx.Errors.String()))
		}

		logger.InfoContext(ctx, options.Msg, attrs...)
	}
}

func roundSize(s int) size.Size {
	if s < 0 {
		return size.NewInt(0, size.B)
	}
	bodysize := size.NewInt(s, size.B)

	if size.Unit(s) > size.GB {
		bodysize = bodysize.To(size.GB)
	} else if size.Unit(s) > size.MB {
		bodysize = bodysize.To(size.MB)
	} else if size.Unit(s) > size.KB {
		bodysize = bodysize.To(size.KB)
	}
	return bodysize
}

func roundDuration(duration time.Duration) string {
	unit := "s"
	base := time.Second

	if duration > time.Minute {
		unit = "min"
		base = time.Minute
	} else if duration > time.Second {
		unit = "ns"
		base = time.Second
	} else if duration > time.Millisecond {
		unit = "ms"
		base = time.Millisecond
	} else if duration > time.Microsecond {
		unit = "µs"
		base = time.Microsecond
	} else if duration > time.Nanosecond {
		unit = "ns"
		base = time.Nanosecond
	}

	return fmt.Sprintf("%.2f%s", float64(duration)/float64(base), unit)
}
