package eventsender

import (
	"context"
	"log/slog"
	"time"
)

type Event interface {
}

type Sender struct {
	Log *slog.Logger
}

// Реализация одного воркера, можно оборачивать в workerpool,
// протягивать context по всей области, но обязательно помнить
// Чтобы доступ к одним и тем же данным не имело > 1 сущности,
// Иначе будет плохо...

func StartWorkerPool(ctx context.Context) {
	// temp, здесь надо сделать воркер пул, в сигнатуру запихать кол-во воркеров про ACID postgres'a помним
}

func (s *Sender) StartProcessEvents(ctx context.Context, handlePeriod time.Duration) {
	const op = "eventsender.StartProcessEvents"

	log := s.Log.With(slog.String("op", op))

	log.Info("start event sending")

	ticker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("event sender is stopped")
				return
			case <-ticker.C:
				// Здесь мы пишем логику обработки
				log.Info("check events")
			}
		}
	}()

}
