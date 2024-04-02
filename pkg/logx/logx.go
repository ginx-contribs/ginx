package logx

import (
	"fmt"
	"github.com/dstgo/filebox"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
)

const (
	TextFormat = "TEXT"
	JSONFormat = "JSON"
)

// Options represents logger configuration
type Options struct {
	Output string `mapstructure:"output"`
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Source bool   `mapstructure:"source"`

	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}

func defaultReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	a.Key = strings.ToUpper(a.Key)
	return a
}

type Option func(options *Options)

func WithLevel(level string) Option {
	return func(options *Options) {
		options.Level = level
	}
}

func WithOutput(output string) Option {
	return func(options *Options) {
		options.Output = output
	}
}

func WithFormat(format string) Option {
	return func(options *Options) {
		options.Format = format
	}
}

func WithSource(source bool) Option {
	return func(options *Options) {
		options.Source = source
	}
}

func WithReplaceAttr(ReplaceAttr func(groups []string, a slog.Attr) slog.Attr) Option {
	return func(options *Options) {
		options.ReplaceAttr = ReplaceAttr
	}
}

// Logger is a simple slog wrapper
type Logger struct {
	l      *slog.Logger
	level  slog.Level
	output *os.File

	writer io.Writer
	opt    Options
}

func (l *Logger) Writer() io.Writer {
	return l.writer
}

func (l *Logger) Slog() *slog.Logger {
	return l.l
}

func (l *Logger) Level() slog.Level {
	return l.level
}

func (l *Logger) Close() error {
	if l.output != nil {
		return l.output.Close()
	}
	return nil
}

// NewTextLogger return a text logger
func NewTextLogger(output string, level string, source bool) (*Logger, error) {
	logger, err := New(
		WithOutput(output),
		WithFormat(TextFormat),
		WithLevel(level),
		WithSource(source),
	)
	return logger, err
}

// NewJsonLogger return a json logger
func NewJsonLogger(output string, level string, source bool) (*Logger, error) {
	logger, err := New(
		WithOutput(output),
		WithFormat(JSONFormat),
		WithLevel(level),
		WithSource(source),
	)
	return logger, err
}

// New return a logger with options
func New(opts ...Option) (*Logger, error) {

	// apply opts
	var opt Options
	for _, option := range opts {
		option(&opt)
	}

	// new logger
	logger := &Logger{opt: opt}

	if opt.Level == "" {
		opt.Level = slog.LevelInfo.String()
	}

	if opt.Format == "" {
		opt.Format = TextFormat
	} else if !slices.Contains([]string{TextFormat, JSONFormat}, opt.Format) {
		return nil, fmt.Errorf("invalid log format: %s", opt.Format)
	}

	var levelvar slog.LevelVar
	if err := levelvar.UnmarshalText([]byte(opt.Level)); err != nil {
		return nil, err
	}
	logger.level = levelvar.Level()

	// output
	if opt.Output == "" {
		logger.writer = os.Stdout
	} else {
		writer, err := filebox.OpenFile(opt.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		multiWriter := io.MultiWriter(os.Stdout, writer)
		logger.writer = multiWriter

		logger.output = writer
	}

	handlerOpt := &slog.HandlerOptions{
		AddSource:   opt.Source,
		Level:       logger.level,
		ReplaceAttr: opt.ReplaceAttr,
	}

	var handler slog.Handler
	if opt.Format == TextFormat {
		handler = slog.NewTextHandler(logger.writer, handlerOpt)
	} else {
		handler = slog.NewJSONHandler(logger.writer, handlerOpt)
	}
	slogger := slog.New(handler)

	logger.l = slogger

	return logger, nil
}
