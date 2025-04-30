package tracing

import (
	"context"
	"payment/internal/telemetry/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func NewTracerProvider(cfg *config.Config) (shutdown func(ctx context.Context) error, err error) {
	// создаем структуры ...opts для каждого пути куда нужны какие-то опшены
	// Потом валидируем через if, составляем через append() структуры
	// Потом создаем что нам надо, ставим как глобальный, и возвращаем функцию для возврата

	var oltpOpts []otlptracehttp.Option        // первым идет для парса
	var traceOpts []trace.TracerProviderOption // Вторым идет

	oltpOpts = append(oltpOpts, otlptracehttp.WithEndpoint(cfg.Traces.Endpoint))

	if cfg.Insecure {
		oltpOpts = append(oltpOpts, otlptracehttp.WithInsecure())
	}

	exp, err := otlptracehttp.New(context.Background(), oltpOpts...)
	if err != nil {
		return nil, err
	}

	traceOpts = append(traceOpts, trace.WithBatcher(exp))
	traceOpts = append(traceOpts, trace.WithResource(resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.ServiceName),
	)))
	traceOpts = append(traceOpts, sampler(cfg.Traces.Sampler, cfg.Traces.SamplerRatio))

	tp := trace.NewTracerProvider(traceOpts...)

	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

// setup sampler for traceOpts
func sampler(sampler string, ratioBased float64) trace.TracerProviderOption {
	switch sampler {
	case "always_on":
		return trace.WithSampler(trace.AlwaysSample())
	case "always_off":
		return trace.WithSampler(trace.NeverSample())
	case "parent_based":
		return trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(ratioBased)))
	default:
		return trace.WithSampler(trace.ParentBased(trace.AlwaysSample()))
	}
}
