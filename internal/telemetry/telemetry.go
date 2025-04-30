package telemetry

import (
	"context"
	"payment/internal/telemetry/config"
	"payment/internal/telemetry/metrics"
	"payment/internal/telemetry/tracing"
	"time"
)

func New(opts ...config.Option) (func(ctx context.Context) error, error) {

	cfg := &config.Config{
		ServiceName: "default-service",
		Insecure:    true,
		Traces: config.TracesConfig{
			Endpoint:     "localhost:4318",
			Timeout:      time.Second * 5,
			Sampler:      "parent_based",
			SamplerRatio: 0.5,
		},
		Metrics: config.MetricsConfig{
			Endpoint: "localhost:4318",
			Timeout:  time.Second * 5,
			Interval: time.Second * 5,
		},
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
	// super mega кривой, but pofig
	fn := func(ctx context.Context) error {
		err := shutdownMetrics(ctx)
		if err != nil {
			return err
		}
		err = shutdownTrace(ctx)
		if err != nil {
			return err
		}
		return nil
	}
	// Возвращаем функцию для выключения (для graceful shutdown)
	return fn, nil
}
