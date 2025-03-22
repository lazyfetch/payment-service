package confirmsrv

import (
	"context"
	"log/slog"
	"payment/internal/domain/models"
)

type PaymentUpdater interface {
	UpdatePayment(ctx context.Context, data *models.DBPayment) error
}

type PaymentProvider interface {
	IsIdempotencyKey(ctx context.Context, data *models.DBPayment) (bool, error)
}

type ConfirmService struct {
	log         *slog.Logger
	paymentupdr PaymentUpdater
}

func New(log *slog.Logger, paymentupdr PaymentUpdater) *ConfirmService {
	return &ConfirmService{
		log:         log,
		paymentupdr: paymentupdr,
	}
}

func (c *ConfirmService) ValidateWebhook() error {

	// валидируем входные данные

	// обращаемся к базе на поиск idempotency_key, если его нету идем дальше

	// потом изменяем статус и Updated_At в DB

	// потом добавляем в очередь кафка новое событие

	return nil // temp
}
