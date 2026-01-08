package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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
	AddCaller   *bool
	AddStack    *bool
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

	zapLoggerOptions := []zap.Option{}

	if config.AddCaller != nil && *config.AddCaller {
		zapLoggerOptions = append(zapLoggerOptions, zap.AddCaller())
	}

	if config.AddStack != nil && *config.AddStack {
		zapLoggerOptions = append(zapLoggerOptions, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	zapLogger := zap.New(core, zapLoggerOptions...)

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

// Sanitize takes a struct and returns a sanitized version with fields marked
// with `log:"redact"` tag having their values replaced with "[REDACTED]".
// The result is returned as a map[string]any suitable for logging.
func Sanitize(input any) (map[string]any, error) {
	if input == nil {
		return nil, nil
	}

	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or pointer to struct, got %T", input)
	}

	result := make(map[string]any)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Skip unexported fields
		if !fieldVal.CanInterface() {
			continue
		}

		// Get the JSON tag name, fallback to field name
		jsonTag := field.Tag.Get("json")
		fieldName := field.Name
		if jsonTag != "" && jsonTag != "-" {
			// Handle json tag options like "omitempty"
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				fieldName = parts[0]
			}
		}

		// Check if field should be redacted
		logTag := field.Tag.Get("log")
		shouldRedact := strings.Contains(logTag, "redact")

		// Handle the field value
		sanitizedValue := sanitizeValue(fieldVal, shouldRedact)
		result[fieldName] = sanitizedValue
	}

	return result, nil
}

// sanitizeValue recursively sanitizes a value based on its type
func sanitizeValue(val reflect.Value, shouldRedact bool) any {
	if shouldRedact {
		return "[REDACTED]"
	}

	if !val.IsValid() {
		return nil
	}

	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}
		return sanitizeValue(val.Elem(), shouldRedact)

	case reflect.Struct:
		result := make(map[string]any)
		typ := val.Type()

		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			fieldVal := val.Field(i)

			if !fieldVal.CanInterface() {
				continue
			}

			jsonTag := field.Tag.Get("json")
			fieldName := field.Name
			if jsonTag != "" && jsonTag != "-" {
				parts := strings.Split(jsonTag, ",")
				if parts[0] != "" {
					fieldName = parts[0]
				}
			}

			logTag := field.Tag.Get("log")
			nestedShouldRedact := strings.Contains(logTag, "redact")

			sanitizedValue := sanitizeValue(fieldVal, nestedShouldRedact)
			result[fieldName] = sanitizedValue
		}

		return result

	case reflect.Slice, reflect.Array:
		if val.IsNil() && val.Kind() == reflect.Slice {
			return nil
		}

		length := val.Len()
		result := make([]any, length)

		for i := 0; i < length; i++ {
			result[i] = sanitizeValue(val.Index(i), shouldRedact)
		}

		return result

	case reflect.Map:
		if val.IsNil() {
			return nil
		}

		result := make(map[string]any)
		for _, key := range val.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			result[keyStr] = sanitizeValue(val.MapIndex(key), shouldRedact)
		}

		return result

	case reflect.Interface:
		if val.IsNil() {
			return nil
		}
		return sanitizeValue(val.Elem(), shouldRedact)

	default:
		// For primitive types, return the value as-is
		return val.Interface()
	}
}

// SanitizeToJSON takes a struct and returns a sanitized JSON string.
// Fields marked with `log:"redact"` will have their values replaced with "[REDACTED]".
func SanitizeToJSON(input any) (string, error) {
	sanitized, err := Sanitize(input)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(sanitized)
	if err != nil {
		return "", fmt.Errorf("failed to marshal sanitized data: %w", err)
	}

	return string(jsonBytes), nil
}
