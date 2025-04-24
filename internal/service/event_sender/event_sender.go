package eventsender

import (
	"context"
	"log/slog"
	"payment/internal/domain/models"
	"payment/internal/lib/logger/sl"
	"time"
)

type EventSender interface {
	SendEvent(ctx context.Context, payload models.Event) error // temp, мб чтото еще надо будет вернуть
}

type EventProvider interface {
	GetEvent(ctx context.Context) (models.Event, error)
}

type Sender struct {
	Log *slog.Logger
	// get info about event
	EventProvider
	// send to broker.. whatever
	EventSender
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

				// Получаем событие через database
				// Обязательно помним про ACID, не берем одно и то же сообщение два раза, избегаем повторок
				// Смотрим на статус!111!
				// Помним об этом, ибо будет воркер пул, который нужно акуратненько трогать.,.,,

				payload, err := s.EventProvider.GetEvent(ctx) // подумать над тем, чтобы правильный контекст протянуть temp
				if err != nil {
					log.Error("got error!", sl.Err(err)) // Через errors.Is() проверяем всю эту шоблу, чтобы корректно логировать
				}
				// temp, мб для логов мы хотим выносить отдельно idempotency key, чтобы трекать полный путь операции, и когда она отправилась
				// так будет лучше понять что пошло не так, офк обязательно добавить openTelemetry
				log.Info("find data! start sending...")
				if err := s.EventSender.SendEvent(ctx, payload); err != nil {
					log.Error("got error!", sl.Err(err)) // Через errors.Is() проверяем всю эту шоблу, чтобы корректно логировать
				}
			}
		}
	}()
}
