package otel

import (
	"context"
	"go-a-b-microservices/pkg/logger"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// InitTracer initializes the OpenTelemetry tracer with Zipkin exporter
func InitTracer(serviceName, zipkinEndpoint string, log logger.Logger) (*sdktrace.TracerProvider, error) {
	// Create Zipkin exporter
	exporter, err := zipkin.New(zipkinEndpoint)
	if err != nil {
		log.Error("Failed to create Zipkin exporter: %v", err)
		return nil, err
	}

	// Create a tracer provider with the Zipkin exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	// Set the global tracer provider and propagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}

// ShutdownTracer gracefully shuts down the tracer
func ShutdownTracer(ctx context.Context, tp *sdktrace.TracerProvider, log logger.Logger) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	if err := tp.Shutdown(ctx); err != nil {
		log.Error("Error shutting down tracer provider: %v", err)
	}
}
