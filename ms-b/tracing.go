package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func initTracer() *trace.TracerProvider {
	ctx := context.Background()
	// exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint("otel-collector:4318"), otlptracehttp.WithInsecure())
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint("otel-collector:4317"), otlptracegrpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("ms-b"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp
}
