# Logger Package

A comprehensive logging package for Go applications with Zap integration, Loki support, and trace context correlation.

## Overview

The Logger package enables you to:
- Use structured logging with Zap
- Integrate with Loki for log aggregation
- Correlate logs with OpenTelemetry trace context
- Sanitize sensitive data before logging
- Configure log levels and output formats
- Support file and console logging
- Use context-aware logging with trace IDs

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/logger
```

## Configuration

### Configuration

Create a logger with configuration:

```go
import (
    "time"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

logger := logger.New(logger.Config{
    Level:       "info",
    JSONOutput:  true,
    FilePath:    "/var/log/app.log",
    Environment: "production",
    Loki: &logger.LokiConfig{
        Enabled:   true,
        URL:       "http://loki:3100/loki/api/v1/push",
        BatchSize: 100,
        BatchWait: 5 * time.Second,
        Labels: map[string]string{
            "service": "my-service",
            "env":     "production",
        },
    },
})
```

### Configuration Fields

#### Config

- **Level** (string, optional): Log level. Options: `"debug"`, `"info"`, `"warn"`, `"error"`, `"fatal"`. Defaults to `"info"`
- **JSONOutput** (bool, optional): Enable JSON output format. Defaults to false (console format)
- **FilePath** (string, optional): File path for log output. If empty, logs to stdout. In development, always logs to stdout
- **Environment** (string, optional): Environment name (e.g., `"development"`, `"production"`). In development, always logs to stdout
- **Loki** (*LokiConfig, optional): Loki configuration for log aggregation

#### LokiConfig

- **Enabled** (bool, required): Enable/disable Loki integration
- **URL** (string, required): Loki push API endpoint (e.g., `"http://loki:3100/loki/api/v1/push"`)
- **BatchSize** (int, optional): Number of entries to batch before sending. Defaults to 100
- **BatchWait** (time.Duration, optional): Maximum time to wait before sending batch. Defaults to 5 seconds
- **Labels** (map[string]string, optional): Default labels to add to all log entries

## Usage

### Basic Logging

```go
import (
    "github.com/TixiaOTA/tixia-pkg/logger"
    "go.uber.org/zap"
)

// Create default logger
logger := logger.Default()

// Log messages
logger.Info("Application started")
logger.Debug("Debug information", zap.String("key", "value"))
logger.Warn("Warning message", zap.Int("count", 42))
logger.Error("Error occurred", zap.Error(err))
```

### Context-Aware Logging with Trace

```go
import (
    "context"
    "github.com/TixiaOTA/tixia-pkg/logger"
    "go.uber.org/zap"
)

ctx := context.Background()
spanID := "span-123"

// Create logger with trace context
traceLogger := logger.WithTrace(ctx, spanID)

// Log with trace context (includes trace_id and span_id)
traceLogger.Info("Processing request", zap.String("user_id", "123"))
```

### Logger with Additional Fields

```go
// Create logger with fields
logger := logger.New(logger.Config{
    Level:      "info",
    JSONOutput: true,
}).With(
    zap.String("service", "api"),
    zap.String("version", "1.0.0"),
)

// All logs will include service and version fields
logger.Info("Request received")
```

### Named Logger

```go
// Create named logger
logger := logger.Default().Named("database")

// All logs will include logger name
logger.Info("Connection established")
```

### Sanitizing Sensitive Data

```go
import "github.com/TixiaOTA/tixia-pkg/logger"

type User struct {
    Username string `json:"username"`
    Password string `json:"password" log:"redact"`
    Email    string `json:"email"`
}

user := User{
    Username: "johndoe",
    Password: "secret123",
    Email:    "john@example.com",
}

// Sanitize before logging
sanitized, err := logger.Sanitize(user)
if err != nil {
    log.Fatal(err)
}

logger.Info("User data", zap.Any("user", sanitized))
// Password field will be "[REDACTED]"

// Or use SanitizeToJSON
jsonStr, err := logger.SanitizeToJSON(user)
logger.Info("User JSON", zap.String("user", jsonStr))
```

### Closing Logger

```go
// Always close logger to flush buffered logs and stop background goroutines
defer logger.Close()
```

## API Reference

### Public Functions

#### New

Creates a new logger with the given configuration.

```go
func New(config Config) *Logger
```

**Parameters:**

- `config` (Config, required): Logger configuration

**Return Values:**

- `*Logger`: Logger instance

#### Default

Creates a default logger with basic configuration.

```go
func Default() *Logger
```

**Return Values:**

- `*Logger`: Default logger instance

**Special Notes:** Uses info level, console output, stdout destination.

#### Sanitize

Takes a struct and returns a sanitized version with fields marked with `log:"redact"` tag having their values replaced with "[REDACTED]".

```go
func Sanitize(input any) (map[string]any, error)
```

**Parameters:**

- `input` (any, required): Struct or pointer to struct to sanitize

**Return Values:**

- `map[string]any`: Sanitized data as map
- `error`: Error if input is not a struct

**Special Notes:** Fields with `log:"redact"` tag will have values replaced with "[REDACTED]". Works recursively for nested structs, slices, and maps.

#### SanitizeToJSON

Takes a struct and returns a sanitized JSON string.

```go
func SanitizeToJSON(input any) (string, error)
```

**Parameters:**

- `input` (any, required): Struct or pointer to struct to sanitize

**Return Values:**

- `string`: Sanitized JSON string
- `error`: Error if sanitization or marshaling fails

### Public Types

#### Logger

Logger struct that wraps zap.Logger with additional features.

```go
type Logger struct {
    *zap.Logger
    lokiClient *loki.Client
    ctx        context.Context
}
```

#### Methods on Logger

##### With

Returns a new Logger with additional fields.

```go
func (l *Logger) With(fields ...zapcore.Field) *Logger
```

**Parameters:**

- `fields` (...zapcore.Field, required): Zap fields to add

**Return Values:**

- `*Logger`: New logger instance with fields

##### WithTrace

Returns a logger with trace context fields.

```go
func (l *Logger) WithTrace(ctx context.Context, spanId string) *Logger
```

**Parameters:**

- `ctx` (context.Context, required): Context with trace span
- `spanId` (string, optional): Span ID. If empty, span_id field will be omitted

**Return Values:**

- `*Logger`: Logger with trace context (trace_id and span_id fields)

**Special Notes:** If spanId is empty, the span_id field will be omitted from log entries.

##### Error

Logs an error message and sets span status if trace context is available.

```go
func (l *Logger) Error(msg string, fields ...zapcore.Field)
```

**Parameters:**

- `msg` (string, required): Error message
- `fields` (...zapcore.Field, optional): Additional fields

**Special Notes:** If logger has trace context, sets span status to error.

##### WithTraceErrorSpan

Sets span status to error and logs error details.

```go
func (l *Logger) WithTraceErrorSpan(ctx context.Context, spanId string, err error)
```

**Parameters:**

- `ctx` (context.Context, required): Context with trace span
- `spanId` (string, required): Span ID
- `err` (error, required): Error to log

**Special Notes:** Sets span status and logs error in one call.

##### Named

Returns a new Logger with the given name.

```go
func (l *Logger) Named(name string) *Logger
```

**Parameters:**

- `name` (string, required): Logger name

**Return Values:**

- `*Logger`: Named logger instance

##### Sync

Flushes any buffered log entries.

```go
func (l *Logger) Sync() error
```

**Return Values:**

- `error`: Error if sync fails

##### Close

Stops the logger and any background goroutines.

```go
func (l *Logger) Close() error
```

**Return Values:**

- `error`: Error if close fails

**Special Notes:** Stops Loki client if enabled and syncs logs.

## Error Handling

The logger handles errors gracefully:

- **Invalid log level**: Defaults to info level
- **File write errors**: Falls back to stdout in development
- **Loki connection errors**: Logs continue to work, Loki errors are logged
- **Sanitization errors**: Returns error if input is not a struct

## Data Models

### Config

Logger configuration:

```go
type Config struct {
    Level       string
    JSONOutput  bool
    FilePath    string
    Environment string
    Loki        *LokiConfig
}
```

### LokiConfig

Loki-specific configuration:

```go
type LokiConfig struct {
    Enabled   bool
    URL       string
    BatchSize int
    BatchWait time.Duration
    Labels    map[string]string
}
```

## Examples

### Complete Example

```go
package main

import (
    "context"
    "time"

    "github.com/TixiaOTA/tixia-pkg/logger"
    "go.uber.org/zap"
)

func main() {
    // Create logger with Loki integration
    logger := logger.New(logger.Config{
        Level:       "info",
        JSONOutput:  true,
        Environment: "production",
        Loki: &logger.LokiConfig{
            Enabled:   true,
            URL:       "http://loki:3100/loki/api/v1/push",
            BatchSize: 100,
            BatchWait: 5 * time.Second,
            Labels: map[string]string{
                "service": "my-service",
            },
        },
    })
    defer logger.Close()

    // Basic logging
    logger.Info("Application started")

    // Logging with fields
    logger.Info("User logged in",
        zap.String("user_id", "123"),
        zap.String("ip", "192.168.1.1"),
    )

    // Context-aware logging
    ctx := context.Background()
    traceLogger := logger.WithTrace(ctx, "span-123")
    traceLogger.Info("Processing request")

    // Error logging
    err := fmt.Errorf("database connection failed")
    logger.Error("Operation failed", zap.Error(err))

    // Sanitized logging
    type Credentials struct {
        Username string `json:"username"`
        Password string `json:"password" log:"redact"`
    }

    creds := Credentials{
        Username: "user",
        Password: "secret",
    }

    sanitized, _ := logger.Sanitize(creds)
    logger.Info("Login attempt", zap.Any("credentials", sanitized))
}
```

### Example with File Logging

```go
logger := logger.New(logger.Config{
    Level:      "info",
    JSONOutput: true,
    FilePath:   "/var/log/app.log",
})

// Logs will be written to file with rotation
logger.Info("Log entry")
```

### Example with Trace Correlation

```go
import (
    "context"
    "github.com/TixiaOTA/tixia-pkg/logger"
    "go.uber.org/zap"
)

func handleRequest(ctx context.Context, logger *logger.Logger, spanID string) {
    // Create logger with trace context
    traceLogger := logger.WithTrace(ctx, spanID)

    // All logs will include trace_id and span_id
    traceLogger.Info("Request started")
    traceLogger.Info("Processing data", zap.String("data", "value"))
    traceLogger.Error("Error occurred", zap.Error(err))
}
```

## Log Sanitization

The logger provides utilities to sanitize sensitive data before logging:

### Using log Tag

Mark fields with `log:"redact"` tag:

```go
type User struct {
    Username string `json:"username"`
    Password string `json:"password" log:"redact"`
    APIKey   string `json:"api_key" log:"redact"`
}
```

### Sanitization Behavior

- Fields with `log:"redact"` tag are replaced with `"[REDACTED]"`
- Works recursively for nested structs
- Handles slices, arrays, and maps
- Preserves JSON tag names
- Skips unexported fields

## License

This package is part of the TixiaOTA project.

