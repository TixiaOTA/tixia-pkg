package logger

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/TixiaOTA/tixia-pkg/logger/loki"
)

// Logger is a wrapper around zap.Logger
type Logger struct {
	*zap.Logger
	lokiClient *loki.Client
	ctx        context.Context
}

// Config represents logger configuration
type Config struct {
	Level       string
	JSONOutput  bool
	FilePath    string
	Environment string
	Loki        *LokiConfig
}

// LokiConfig represents Loki-specific configuration
type LokiConfig struct {
	Enabled   bool
	URL       string
	BatchSize int
	BatchWait time.Duration
	Labels    map[string]string
}

// New creates a new logger with the given configuration
func New(config Config) *Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if config.JSONOutput {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var core zapcore.Core
	var lokiClient *loki.Client

	cores := []zapcore.Core{}

	// In development environment, always log to stdout
	if config.Environment == "development" {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), parseLevel(config.Level)))
	} else if config.FilePath != "" {
		writer := &lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    100, // MB
			MaxBackups: 5,
			MaxAge:     30, // days
			Compress:   true,
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(writer), parseLevel(config.Level)))
	} else {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), parseLevel(config.Level)))
	}

	if config.Loki != nil && config.Loki.Enabled && config.Loki.URL != "" {
		lokiClient = loki.NewClient(loki.Config{
			URL:       config.Loki.URL,
			BatchSize: config.Loki.BatchSize,
			BatchWait: config.Loki.BatchWait,
			Labels:    config.Loki.Labels,
		})

		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(&lokiWriter{client: lokiClient}),
			parseLevel(config.Level),
		))
	}

	core = zapcore.NewTee(cores...)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		Logger:     zapLogger,
		lokiClient: lokiClient,
	}
}

// Default creates a default logger
func Default() *Logger {
	return New(Config{
		Level:      "info",
		JSONOutput: false,
		FilePath:   "",
	})
}

// With returns a new Logger with additional fields
func (l *Logger) With(fields ...zapcore.Field) *Logger {
	return &Logger{
		Logger:     l.Logger.With(fields...),
		lokiClient: l.lokiClient,
	}
}

// WithTrace returns a logger with trace context fields.
// If spanId is empty, the span_id field will be omitted from the log entry.
func (l *Logger) WithTrace(ctx context.Context, spanId string) *Logger {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return l
	}

	fields := make([]zap.Field, 0, 2) // Pre-allocate for 2 fields
	fields = append(fields, zap.String("trace_id", span.SpanContext().TraceID().String()))

	if spanId != "" {
		fields = append(fields, zap.String("span_id", spanId))
	}

	return &Logger{
		Logger:     l.Logger.With(fields...),
		lokiClient: l.lokiClient,
		ctx:        ctx, // Store the context
	}
}

func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	if ce := l.Check(zapcore.ErrorLevel, msg); ce != nil {
		ce.Write(fields...)

		// Derive span from the stored context
		span := trace.SpanFromContext(l.ctx)
		if span.SpanContext().IsValid() {
			span.SetStatus(codes.Error, msg)
		}
	}
}

func (l *Logger) WithTraceErrorSpan(ctx context.Context, spanId string, err error) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return
	}

	span.SetStatus(codes.Error, err.Error())
	l.WithTrace(ctx, spanId).Error("error_details", zap.String("error", err.Error()))
}

// Named returns a new Logger with the given name
func (l *Logger) Named(name string) *Logger {
	return &Logger{
		Logger:     l.Logger.Named(name),
		lokiClient: l.lokiClient,
	}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Close stops the logger and any background goroutines
func (l *Logger) Close() error {
	if l.lokiClient != nil {
		l.lokiClient.Stop()
	}
	return l.Sync()
}

// Helper function to parse log level
func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// lokiWriter implements zapcore.WriteSyncer for Loki
type lokiWriter struct {
	client *loki.Client
}

func (w *lokiWriter) Write(p []byte) (n int, err error) {
	level := "info"
	if len(p) > 0 {
		lowered := string(p)
		switch {
		case contains(lowered, "debug"):
			level = "debug"
		case contains(lowered, "info"):
			level = "info"
		case contains(lowered, "warn"):
			level = "warn"
		case contains(lowered, "error"):
			level = "error"
		case contains(lowered, "fatal"):
			level = "fatal"
		}
	}

	w.client.Log(time.Now().UTC(), level, string(p))
	return len(p), nil
}

func (w *lokiWriter) Sync() error {
	return nil
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
