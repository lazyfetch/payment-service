package tracing

import (
	"context"
	"payment/internal/telemetry/config"
)

func NewTracerProvider(cfg *config.Config) (shutdown func(ctx context.Context) error, err error) {
	// Здесь создаешь сам TracerProvider на основе конфига
	return nil, nil
}
