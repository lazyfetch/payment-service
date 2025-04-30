package metrics

import (
	"context"
	"payment/internal/telemetry/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

func NewMetricsProvider(cfg *config.Config) (shutdown func(ctx context.Context) error, err error) {
	var opts []otlpmetrichttp.Option

	opts = append(opts, otlpmetrichttp.WithEndpoint(cfg.Metrics.Endpoint))
	if cfg.Insecure {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}
	opts = append(opts, otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression))
	opts = append(opts, otlpmetrichttp.WithURLPath("/v1/metrics"))

	exp, err := otlpmetrichttp.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	reader := metric.NewPeriodicReader(exp, metric.WithInterval(cfg.Metrics.Interval))

	mp := metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
		)),
	)

	otel.SetMeterProvider(mp)

	return mp.Shutdown, nil
}
