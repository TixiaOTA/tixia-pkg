package otelpkg

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func StartSpanWithID(ctx context.Context, layer, funcName, spanName string) (context.Context, trace.Span, string) {
	tracer := otel.Tracer(fmt.Sprintf("%s/%s", layer, funcName))
	ctx, span := tracer.Start(ctx, spanName)
	return ctx, span, span.SpanContext().SpanID().String()
}

func InjectTraceContext(ctx context.Context) map[string]string {
	carrier := make(map[string]string)
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(carrier))
	return carrier
}

func ExtractTraceContext(ctx context.Context, traceMap map[string]string) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(traceMap))
}

func LogTraceInfo(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()
	if !sc.IsValid() {
		fmt.Println("No valid span context found")
		return
	}
	fmt.Printf("TraceID: %s, SpanID: %s, TraceFlags: %s\n",
		sc.TraceID().String(),
		sc.SpanID().String(),
		sc.TraceFlags().String(),
	)
}
