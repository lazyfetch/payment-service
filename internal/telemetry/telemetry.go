package telemetry

import (
	"context"
	"payment/internal/telemetry/config"
	"payment/internal/telemetry/metrics"
	"payment/internal/telemetry/tracing"
)

func New(opts ...config.Option) (func(ctx context.Context), error) {

	cfg := &config.Config{
		Endpoint: "localhost:4317",
		Service:  "default",
		Insecure: false,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	shutdownTrace, err := tracing.NewTracerProvider(cfg)
	if err != nil {

	}
	shutdownMetrics, err := metrics.NewMetricsProvider(cfg)
	if err != nil {

	}

	// Возвращаем функцию для выключения (для graceful shutdown)
	return func(ctx context.Context) {
		shutdownTrace(ctx) // temp, error handl
		shutdownMetrics(ctx)
	}, nil
}
