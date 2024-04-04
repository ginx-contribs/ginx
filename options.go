package ginx

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

// HookFn will be executed at specific point
type HookFn = func(ctx context.Context) error

type Option func(options *Server)

type TLSOptions struct {
	// TLS Key file
	Key string `mapstructure:"key"`
	// TLS Certificate file
	Cert string `mapstructure:"cert"`
}

type Options struct {
	Mode string `mapstructure:"mode"`
	// specifies the TCP address for the server to listen on,
	// in the form "host:port". If empty, ":http" (port 80) is used
	Address string `mapstructure:"address"`

	// prefix in log records
	LogPrefix string `mapstructure:"logPrefix"`

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. A zero or negative value means
	// there will be no timeout.
	ReadTimeout time.Duration `mapstructure:"readTimeout"`

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers.
	ReadHeaderTimeout time.Duration `mapstructure:"readHeaderTimeout"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response.
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled.
	IdleTimeout time.Duration `mapstructure:"idleTimeout"`

	// MaxMultipartMemory value of 'maxMemory' param that is given to http.Request's ParseMultipartForm
	// method call.
	MaxMultipartMemory int64 `mapstructure:"xaxMultipartMemory"`

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line.
	MaxHeaderBytes int `mapstructure:"maxHeaderBytes"`

	// simple TLS config
	TLS *TLSOptions `mapstructure:"tls"`

	// max wait time after server shutdown
	MaxShutdownTimeout time.Duration `mapstructure:"maxShutdownTimeout"`
}

// WithCtx apply the server root context
func WithCtx(ctx context.Context) Option {
	return func(server *Server) {
		server.ctx = ctx
	}
}

func WithOnShutdown(hooks ...HookFn) Option {
	return func(server *Server) {
		server.OnShutdown = append(server.OnShutdown, hooks...)
	}
}

func WithBeforeStarting(hooks ...HookFn) Option {
	return func(server *Server) {
		server.BeforeStarting = append(server.BeforeStarting, hooks...)
	}
}

func WithAfterStarted(hooks ...HookFn) Option {
	return func(server *Server) {
		server.AfterStarted = append(server.AfterStarted, hooks...)
	}
}

func WithStopSignals(signals ...os.Signal) Option {
	return func(server *Server) {
		server.stopSignals = signals
	}
}

// WithEngine apply a custom engine
func WithEngine(engine *gin.Engine) Option {
	return func(server *Server) {
		server.engine = engine
	}
}

// WithHttpServer apply a custom httpserver
func WithHttpServer(httpserver *http.Server) Option {
	return func(server *Server) {
		server.httpserver = httpserver
	}
}

func WithMiddlewares(handlers ...gin.HandlerFunc) Option {
	return func(server *Server) {
		server.middlewares = append(server.middlewares, handlers...)
	}
}

func WithNoRoute(handlers ...gin.HandlerFunc) Option {
	return func(server *Server) {
		server.noRoute = append(server.noRoute, handlers...)
	}
}

func WithNoMethod(handlers ...gin.HandlerFunc) Option {
	return func(server *Server) {
		server.noMethod = append(server.noMethod, handlers...)
	}
}

// WithOptions apply a whole options
func WithOptions(options Options) Option {
	return func(server *Server) {
		server.options = options
	}
}

func WithLogPrefix(prefix string) Option {
	return func(server *Server) {
		server.options.LogPrefix = prefix
	}
}

func WithAddress(address string) Option {
	return func(server *Server) {
		server.options.Address = address
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(server *Server) {
		server.options.ReadTimeout = timeout
	}
}

func WithReadHeaderTimeout(timeout time.Duration) Option {
	return func(server *Server) {
		server.options.ReadHeaderTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(server *Server) {
		server.options.WriteTimeout = timeout
	}
}

func WithIdleTimeout(timeout time.Duration) Option {
	return func(server *Server) {
		server.options.IdleTimeout = timeout
	}
}
func WithMultipartMem(mem int64) Option {
	return func(server *Server) {
		server.options.MaxMultipartMemory = mem
	}
}

func WithMaxHeaderBytes(bytes int) Option {
	return func(server *Server) {
		server.options.MaxHeaderBytes = bytes
	}
}

func WithMaxShutdownWait(timeout time.Duration) Option {
	return func(server *Server) {
		server.options.MaxShutdownTimeout = timeout
	}
}

func WithMode(mode string) Option {
	return func(server *Server) {
		server.options.Mode = mode
	}
}

func WithTLS(key string, cert string) Option {
	return func(server *Server) {
		server.options.TLS = &TLSOptions{Key: key, Cert: cert}
	}
}
