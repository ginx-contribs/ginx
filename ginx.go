package ginx

import (
	"context"
	"errors"
	"fmt"
	"github.com/dstgo/size"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx/middleware"
	cmap "github.com/orcaman/concurrent-map/v2"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// GinSilence decide whether to keep gin.DefaultWriter and gin.DefaultErrorWriter silence
var GinSilence = true

func init() {
	if GinSilence {
		gin.DefaultWriter = io.Discard
		gin.DefaultErrorWriter = io.Discard
	}
}

// Default returns *Server with default middlewares
func Default() *Server {
	return New(
		WithNoRoute(middleware.NoRoute()),
		WithNoMethod(middleware.NoMethod(allowMethods...)),
		WithMiddlewares(
			middleware.Logger(slog.Default(), logPrefix),
			middleware.Recovery(slog.Default(), nil),
		),
	)
}

// New returns a new server instance
func New(options ...Option) *Server {

	server := new(Server)
	server.metadata = &FrozenMap[string, routeMeta]{
		m: cmap.New[routeMeta](),
	}

	for _, option := range options {
		option(server)
	}

	if server.ctx == nil {
		server.ctx = context.Background()
	}

	if server.options.Mode == "" {
		server.options.Mode = gin.ReleaseMode
	}
	gin.SetMode(server.options.Mode)

	if len(server.stopSignals) == 0 {
		server.stopSignals = []os.Signal{syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT}
	}

	if server.options.Address == "" {
		server.options.Address = ":8080"
	}

	if server.options.MaxShutdownTimeout == 0 {
		server.options.MaxShutdownTimeout = time.Second * 5
	}

	if server.options.ReadTimeout == 0 {
		server.options.ReadTimeout = time.Second * 60
	}

	if server.options.ReadHeaderTimeout == 0 {
		server.options.ReadHeaderTimeout = time.Second * 60
	}

	if server.options.WriteTimeout == 0 {
		server.options.WriteTimeout = time.Second * 60
	}

	if server.options.IdleTimeout == 0 {
		server.options.IdleTimeout = time.Minute * 5
	}

	if server.options.MaxMultipartMemory == 0 {
		server.options.MaxMultipartMemory = int64(size.MB * 10)
	}

	if server.options.MaxHeaderBytes == 0 {
		server.options.MaxHeaderBytes = int(size.MB)
	}

	server.applyOptions()

	return server
}

// Server is a simple wrapper for http.Server and *gin.Engine, which is more convenient to use.
// It provides hooks can be executed at certain time, ability to graceful shutdown.
type Server struct {
	ctx context.Context

	// truly running server
	httpserver *http.Server

	engine *gin.Engine

	noRoute  gin.HandlersChain
	noMethod gin.HandlersChain
	// global middlewares
	middlewares gin.HandlersChain

	// hooks func
	BeforeStarting []HookFn
	AfterStarted   []HookFn
	OnShutdown     []HookFn

	// metadata is a ready-only map during server running, which holds all the route metadata.
	// It is not thread-safe, should not be modified after server running.
	metadata *FrozenMap[string, routeMeta]

	// os stop signals
	stopSignals []os.Signal

	options Options
}

func (s *Server) HttpServer() *http.Server {
	return s.httpserver
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

// Run just run the http server without hooks, you should use *Server.Spin` im most time.
func (s *Server) Run() error {
	s.metadata.Frozen()
	if s.options.TLS != nil {
		slog.InfoContext(s.ctx, "tls certificate has been configured")
		slog.InfoContext(s.ctx, fmt.Sprintf("server is listiening at %v", s.options.Address))
		return s.httpserver.ListenAndServeTLS(s.options.TLS.Cert, s.options.TLS.Key)
	} else {
		slog.InfoContext(s.ctx, fmt.Sprintf("server is listiening at %v", s.options.Address))
		return s.httpserver.ListenAndServe()
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpserver.Shutdown(ctx)
}

// Spin runs the server in another go routine, and listening for os signals to graceful shutdown,
func (s *Server) Spin() error {
	notifyContext, signalCancel := signal.NotifyContext(s.ctx, s.stopSignals...)
	defer signalCancel()

	slog.DebugContext(s.ctx, "hooks before starting are executing")
	// execute before starting hooks
	err := s.executeHooks(notifyContext, s.BeforeStarting...)
	if err != nil {
		return err
	}

	runCh := make(chan error)

	go func() {
		runCh <- s.Run()
		close(runCh)
	}()

	slog.DebugContext(s.ctx, "hooks after starting are executing")
	// execute after started hooks
	err = s.executeHooks(notifyContext, s.AfterStarted...)
	if err != nil {
		return err
	}

	// wait for server closed or stop signal
	select {
	case <-notifyContext.Done():
		slog.InfoContext(s.ctx, fmt.Sprintf("received stop signal, it will shutdown in %s at latest", s.options.MaxShutdownTimeout.String()))
	case err := <-runCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(s.ctx, "running failed", slog.Any("error", err))
		} else {
			slog.InfoContext(s.ctx, "server closed")
		}
	}

	// ready to server shutdown
	shutdownCh := make(chan error)
	timeoutCtx, shutdownCancel := context.WithTimeout(s.ctx, s.options.MaxShutdownTimeout)
	defer shutdownCancel()

	_ = s.Shutdown(timeoutCtx)

	go func() {
		slog.DebugContext(s.ctx, "hooks on shutdown are executing")
		shutdownCh <- s.executeHooks(timeoutCtx, s.OnShutdown...)
		close(shutdownCh)
	}()

	// wait timeout for execute shutdown hooks
	select {
	case <-timeoutCtx.Done():
		slog.ErrorContext(s.ctx, "shutdown timeout")
	case err := <-shutdownCh:
		if err != nil {
			slog.ErrorContext(s.ctx, "shutdown error", slog.Any("error", err))
			return err
		} else {
			slog.InfoContext(s.ctx, "server shutdown")
		}
	}

	// server running finished
	return nil
}

func (s *Server) executeHooks(ctx context.Context, hooks ...HookFn) error {
	for _, hook := range hooks {
		if err := hook(ctx); err != nil {
			return err
		}
	}
	return nil
}

// applyOptions applies options to http server and engine
func (s *Server) applyOptions() {
	if s.httpserver == nil {
		s.httpserver = &http.Server{}
	}
	if s.engine == nil {
		s.engine = gin.New()
	}

	if s.httpserver.Addr == "" {
		s.httpserver.Addr = s.options.Address
	}

	if s.httpserver.ReadTimeout == 0 {
		s.httpserver.ReadTimeout = s.options.ReadTimeout
	}

	if s.httpserver.ReadHeaderTimeout == 0 {
		s.httpserver.ReadHeaderTimeout = s.options.ReadHeaderTimeout
	}

	if s.httpserver.WriteTimeout == 0 {
		s.httpserver.WriteTimeout = s.options.WriteTimeout
	}

	if s.httpserver.MaxHeaderBytes == 0 {
		s.httpserver.MaxHeaderBytes = s.options.MaxHeaderBytes
	}

	if s.httpserver.Handler != nil {
		if engine, ok := s.httpserver.Handler.(*gin.Engine); ok {
			// overlay engine
			s.engine = engine
		} else {
			panic(fmt.Errorf("expected: github.com/gin-gonic/*gin.Engine, but got %T", s.httpserver.Handler))
		}
	} else {
		// use engine for httpserver handler
		s.httpserver.Handler = s.engine
	}

	if s.engine.MaxMultipartMemory == 0 {
		s.engine.MaxMultipartMemory = s.options.MaxMultipartMemory
	}

	// apply middlewares
	s.engine.Use(metaDataHandler(s.metadata))
	s.engine.Use(s.middlewares...)
	s.engine.NoMethod(s.noMethod...)
	s.engine.NoRoute(s.noRoute...)
}
