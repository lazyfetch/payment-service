package metrics

import (
	"context"
	"payment/internal/telemetry/config"
)

func NewMetricsProvider(cfg *config.Config) (shutdown func(ctx context.Context) error, err error) {
	return
}
