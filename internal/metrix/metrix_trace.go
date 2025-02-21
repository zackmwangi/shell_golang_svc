package metrix

import (
	"github.com/zackmwangi/shell_golang_svc/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func InitTraceProvider(c *config.AppConfig) (*sdktrace.TracerProvider, error) {
	//

	traceExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(c.OtelExporterOtlpURL)))
	//
	if err != nil {
		return nil, err
	}
	//
	if err != nil {
		return nil, err
	}
	//
	metricExporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}
	//
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(metricExporter),
	)
	//
	//serviceResource
	serviceResource, rErr := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			attribute.String("environment", c.AppEnv),
			attribute.String("service.name", c.SvcName),
			attribute.String("service.version", c.SvcVersion),
		),
	)

	if rErr != nil {
		return nil, rErr
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(serviceResource),
	)

	//
	otel.SetMeterProvider(mp)
	otel.SetTracerProvider(tp)
	//
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
	//
}
