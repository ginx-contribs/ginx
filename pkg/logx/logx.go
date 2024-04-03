package logx

import (
	"fmt"
	"github.com/dstgo/filebox"
	"github.com/lmittmann/tint"
	"io"
	"log/slog"
	"os"
	"slices"
	"time"
)

const (
	TextFormat = "TEXT"
	JSONFormat = "JSON"
)

// Options represents logger configuration
type Options struct {
	Output string `mapstructure:"output"`
	// Log level, default INFO
	Level string `mapstructure:"level"`
	// TEXT or JSON
	Format string `mapstructure:"format"`
	// whether to show source files
	Source bool `mapstructure:"source"`
	// custom time format
	TimeFormat string `mapstructure:"timeFormat"`
	// color log only available in TEXT format
	NoColor bool `mapstructure:"color"`
	// attributes replace func
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
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

func WithSource() Option {
	return func(options *Options) {
		options.Source = true
	}
}

func WithTimeFormat(format string) Option {
	return func(options *Options) {
		options.TimeFormat = format
	}
}

func WithColor() Option {
	return func(options *Options) {
		options.NoColor = false
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

	if opt.TimeFormat == "" {
		opt.TimeFormat = time.DateTime
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
		AddSource: opt.Source,
		Level:     logger.level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				a.Value = slog.AnyValue(a.Value.Time().Format(opt.TimeFormat))
			}
			if opt.ReplaceAttr != nil {
				return opt.ReplaceAttr(groups, a)
			}
			return a
		},
	}

	var handler slog.Handler
	if opt.Format == TextFormat {
		handler = tint.NewHandler(logger.writer, &tint.Options{
			AddSource:   opt.Source,
			Level:       logger.level,
			ReplaceAttr: handlerOpt.ReplaceAttr,
			TimeFormat:  opt.TimeFormat,
			NoColor:     opt.NoColor,
		})
	} else {
		handler = slog.NewJSONHandler(logger.writer, handlerOpt)
	}
	slogger := slog.New(handler)

	logger.l = slogger

	return logger, nil
}
