package eventsender

import (
	"context"
	"log/slog"
	"time"
)

type Sender struct {
	log *slog.Logger
}

func (s *Sender) StartProcessEvents(ctx context.Context, handlePeriod time.Duration) {
	const op = "event-sender.StartProcessEvents"

	// log := s.log.With(slog.String("op", op))

	ticker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				//
			}
		}
	}()

}
