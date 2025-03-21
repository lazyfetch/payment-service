package confirmsrv

import (
	"context"
	"log/slog"
	"payment/internal/domain/models"
)

type PaymentUpdater interface {
	UpdatePayment(ctx context.Context, data models.DBPayment) error
}

type PaymentProvider interface {
	Payment() models.GRPCPayment
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

	// сначала используем валидатор robokassa

	// потом изменяем статус и Updated_At в DB

	// потом добавляем в очередь кафка новое событие

	return nil // temp
}
