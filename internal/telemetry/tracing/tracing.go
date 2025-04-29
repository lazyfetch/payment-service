package tracing

import "context"

func NewTracerProvider() (shutdown func(ctx context.Context) error, err error) {
	// Здесь создаешь сам TracerProvider на основе конфига
	return nil, nil
}
