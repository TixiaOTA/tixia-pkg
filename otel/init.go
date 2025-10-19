package otelpkg

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.15.0"
	"go.uber.org/zap"
)

func InitTracer(tempoURL string, logger *zap.Logger, serviceName string) func() {
	// Configure OTLP over HTTP exporter to Tempo
	ctx := context.Background()

	// Enable debug logging for OpenTelemetry
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.Error("OpenTelemetry error", zap.Error(err))
	}))

	httpClient := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(tempoURL),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         true,
			InitialInterval: 5 * time.Second,
			MaxInterval:     30 * time.Second,
			MaxElapsedTime:  time.Minute,
		}),
	)

	exp, err := otlptrace.New(ctx, httpClient)
	if err != nil {
		logger.Fatal("failed to create OTLP exporter", zap.Error(err))
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		logger.Fatal("failed to create resource", zap.Error(err))
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return func() { _ = tp.Shutdown(ctx) }
}
