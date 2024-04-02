package requestid

import (
	"github.com/246859/ginx/constant/headers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Options struct {
	Header    string
	Generator func(ctx *gin.Context) string
}

type Option func(options *Options)

func WithHeader(header string) Option {
	return func(options *Options) {
		options.Header = header
	}
}

func WithGenerator(gen func(ctx *gin.Context) string) Option {
	return func(options *Options) {
		options.Generator = gen
	}
}

// RequestId returns the request-id gin handler
func RequestId(opts ...Option) gin.HandlerFunc {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	if options.Header == "" {
		options.Header = headers.XRequestId
	}

	if options.Generator == nil {
		options.Generator = func(_ *gin.Context) string {
			return uuid.NewString()
		}
	}

	return func(ctx *gin.Context) {
		// if client has already carried with request-id
		id := ctx.Request.Header.Get(options.Header)
		if id == "" { // generate a new one
			id = options.Generator(ctx)
		}
		// set in response
		ctx.Writer.Header().Set(options.Header, id)
	}
}
