# OpenTelemetry Package

A utility package for OpenTelemetry distributed tracing integration with Tempo, providing span creation and trace context management.

## Overview

The OpenTelemetry package enables you to:
- Initialize OpenTelemetry tracer with Tempo integration
- Create spans with automatic span ID generation
- Manage trace context propagation
- Inject and extract trace context for distributed systems
- Correlate logs with trace context

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/otel
```

## Configuration

### Configuration

Initialize the OpenTelemetry tracer:

```go
import (
    "github.com/TixiaOTA/tixia-pkg/otel"
    "go.uber.org/zap"
)

logger, _ := zap.NewProduction()

// Initialize tracer with Tempo
shutdown := otel.InitTracer(
    "http://tempo:4318", // Tempo OTLP endpoint
    logger,
    "my-service",        // Service name
)

// Always call shutdown when application exits
defer shutdown()
```

### Configuration Fields

#### InitTracer Parameters

- **tempoURL** (string, required): Tempo OTLP HTTP endpoint (e.g., `"http://tempo:4318"`)
- **logger** (*zap.Logger, required): Logger instance for OpenTelemetry errors
- **serviceName** (string, required): Service name for trace identification

**Return Values:**

- `func()`: Shutdown function to call when application exits

## Usage

### Basic Span Creation

```go
import (
    "context"
    "github.com/TixiaOTA/tixia-pkg/otel"
)

func processRequest(ctx context.Context) {
    // Start a span
    ctx, span, spanID := otel.StartSpanWithID(ctx, "provider", "mg", "SearchHotel")
    defer span.End()

    // Use spanID for logging correlation
    logger.WithTrace(ctx, spanID).Info("Processing request")

    // Your business logic here
    // ...

    // Span ends automatically when defer executes
}
```

### Trace Context Propagation

```go
import "github.com/TixiaOTA/tixia-pkg/otel"

// Inject trace context into headers for HTTP requests
func makeHTTPRequest(ctx context.Context) {
    traceHeaders := otel.InjectTraceContext(ctx)
    
    // traceHeaders contains trace context as map[string]string
    // Use these headers in your HTTP request
    for key, value := range traceHeaders {
        req.Header.Set(key, value)
    }
}

// Extract trace context from incoming HTTP request
func handleHTTPRequest(req *http.Request) context.Context {
    traceMap := make(map[string]string)
    traceMap["traceparent"] = req.Header.Get("traceparent")
    traceMap["tracestate"] = req.Header.Get("tracestate")
    
    ctx := otel.ExtractTraceContext(context.Background(), traceMap)
    return ctx
}
```

### Logging Trace Information

```go
import "github.com/TixiaOTA/tixia-pkg/otel"

// Log trace information for debugging
func debugTrace(ctx context.Context) {
    otel.LogTraceInfo(ctx)
    // Output: TraceID: abc123..., SpanID: def456..., TraceFlags: 01
}
```

## API Reference

### Public Functions

#### InitTracer

Initializes OpenTelemetry tracer with Tempo integration.

```go
func InitTracer(tempoURL string, logger *zap.Logger, serviceName string) func()
```

**Parameters:**

- `tempoURL` (string, required): Tempo OTLP HTTP endpoint
- `logger` (*zap.Logger, required): Logger for OpenTelemetry errors
- `serviceName` (string, required): Service name for traces

**Return Values:**

- `func()`: Shutdown function to call on application exit

**Special Notes:** Configures OTLP exporter with retry logic. Sets up trace context and baggage propagation. Always call the returned shutdown function to flush traces.

#### StartSpanWithID

Starts a new span and returns context, span, and span ID.

```go
func StartSpanWithID(ctx context.Context, layer, funcName, spanName string) (context.Context, trace.Span, string)
```

**Parameters:**

- `ctx` (context.Context, required): Parent context
- `layer` (string, required): Layer name (e.g., `"provider"`, `"usecase"`)
- `funcName` (string, required): Function/package name (e.g., `"mg"`, `"vervotech"`)
- `spanName` (string, required): Span name (e.g., `"SearchHotel"`, `"BookHotel"`)

**Return Values:**

- `context.Context`: Context with span
- `trace.Span`: OpenTelemetry span
- `string`: Span ID as string

**Special Notes:** Tracer name is formatted as `"{layer}/{funcName}"`. Span ID can be used for log correlation.

#### InjectTraceContext

Injects trace context into a map for propagation.

```go
func InjectTraceContext(ctx context.Context) map[string]string
```

**Parameters:**

- `ctx` (context.Context, required): Context with trace span

**Return Values:**

- `map[string]string`: Map containing trace context headers

**Special Notes:** Returns headers like `traceparent` and `tracestate` for HTTP propagation.

#### ExtractTraceContext

Extracts trace context from a map and creates a new context.

```go
func ExtractTraceContext(ctx context.Context, traceMap map[string]string) context.Context
```

**Parameters:**

- `ctx` (context.Context, required): Base context
- `traceMap` (map[string]string, required): Map containing trace context headers

**Return Values:**

- `context.Context`: Context with extracted trace context

**Special Notes:** Used to continue traces from incoming HTTP requests.

#### LogTraceInfo

Logs trace information for debugging.

```go
func LogTraceInfo(ctx context.Context)
```

**Parameters:**

- `ctx` (context.Context, required): Context with trace span

**Special Notes:** Prints trace ID, span ID, and trace flags to stdout. Useful for debugging.

## Error Handling

The package handles errors gracefully:

- **Tempo connection errors**: Logged but don't prevent application startup
- **Invalid trace context**: Extraction returns context without trace if invalid
- **Missing span context**: Functions handle missing spans gracefully

## Data Models

The package uses OpenTelemetry standard types:

- `trace.Span`: OpenTelemetry span interface
- `trace.SpanContext`: Span context with trace ID and span ID

## Examples

### Complete Example

```go
package main

import (
    "context"
    "net/http"
    "time"

    "github.com/TixiaOTA/tixia-pkg/otel"
    "go.uber.org/zap"
)

func main() {
    // Initialize logger
    logger, _ := zap.NewProduction()

    // Initialize tracer
    shutdown := otel.InitTracer(
        "http://tempo:4318",
        logger,
        "my-service",
    )
    defer shutdown()

    // Example: HTTP handler with tracing
    http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        // Extract trace context from request
        traceMap := make(map[string]string)
        traceMap["traceparent"] = r.Header.Get("traceparent")
        traceMap["tracestate"] = r.Header.Get("tracestate")

        ctx := otel.ExtractTraceContext(r.Context(), traceMap)

        // Start span
        ctx, span, spanID := otel.StartSpanWithID(ctx, "api", "users", "GetUsers")
        defer span.End()

        // Use spanID for logging
        logger.Info("Handling request",
            zap.String("span_id", spanID),
        )

        // Your business logic
        // ...

        w.WriteHeader(http.StatusOK)
    })

    http.ListenAndServe(":8080", nil)
}
```

### Example with Provider Package

```go
import (
    "context"
    "github.com/TixiaOTA/tixia-pkg/otel"
    "github.com/TixiaOTA/tixia-pkg/providers/mg"
)

func searchHotels(ctx context.Context, client *mg.Client) {
    // Start span following provider pattern
    ctx, span, spanID := otel.StartSpanWithID(ctx, "provider", "mg", "SearchHotel")
    defer span.End()

    // Use spanID for logging correlation
    logger.WithTrace(ctx, spanID).Info("Starting hotel search")

    // Make API call
    req := mg.SearchHotelRequest{...}
    resp, err := client.SearchHotel(ctx, req)

    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return
    }

    logger.WithTrace(ctx, spanID).Info("Search completed",
        zap.Int("hotels_found", resp.NoOfHotels),
    )
}
```

### Example with Trace Propagation

```go
import (
    "context"
    "net/http"
    "github.com/TixiaOTA/tixia-pkg/otel"
)

// Client side: Inject trace context
func makeExternalAPICall(ctx context.Context) {
    // Start span
    ctx, span, _ := otel.StartSpanWithID(ctx, "client", "api", "ExternalCall")
    defer span.End()

    // Inject trace context
    traceHeaders := otel.InjectTraceContext(ctx)

    // Create HTTP request
    req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)

    // Add trace headers
    for key, value := range traceHeaders {
        req.Header.Set(key, value)
    }

    // Make request
    http.DefaultClient.Do(req)
}

// Server side: Extract trace context
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Extract trace context
    traceMap := make(map[string]string)
    traceMap["traceparent"] = r.Header.Get("traceparent")
    traceMap["tracestate"] = r.Header.Get("tracestate")

    ctx := otel.ExtractTraceContext(r.Context(), traceMap)

    // Continue trace
    ctx, span, _ := otel.StartSpanWithID(ctx, "server", "handler", "ProcessRequest")
    defer span.End()

    // Process request
    // ...
}
```

## Integration with Logger

The package works seamlessly with the logger package:

```go
import (
    "context"
    "github.com/TixiaOTA/tixia-pkg/otel"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

func process(ctx context.Context, logger *logger.Logger) {
    // Start span
    ctx, span, spanID := otel.StartSpanWithID(ctx, "service", "process", "ProcessData")
    defer span.End()

    // Use spanID for log correlation
    traceLogger := logger.WithTrace(ctx, spanID)
    traceLogger.Info("Processing started")

    // All logs will include trace_id and span_id
    traceLogger.Info("Step completed", zap.String("step", "validation"))
}
```

## License

This package is part of the TixiaOTA project.

