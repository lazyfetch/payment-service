package confirmsrv

import (
	"context"
	"log/slog"
	"payment/internal/govnokassa"
	"payment/internal/lib/logger/sl"
	"time"
)

type PaymentUpdater interface {
	UpdatePayment(ctx context.Context, idemKey string, status string, updatedAt time.Time) error
}

type PaymentProvider interface {
	IdempotencyAndStatus(ctx context.Context, idempotencyKey string) (bool, error)
}

type Validate interface {
	ValidateData(rawData []byte) (*govnokassa.GovnoPayment, error)
}

type ConfirmService struct {
	log         *slog.Logger
	paymentupdr PaymentUpdater
	paymentprv  PaymentProvider
	validate    Validate
}

func New(log *slog.Logger, paymentupdr PaymentUpdater, paymentprv PaymentProvider, validate Validate) *ConfirmService {
	return &ConfirmService{
		log:         log,
		paymentupdr: paymentupdr,
		paymentprv:  paymentprv,
		validate:    validate,
	}
}

func (c *ConfirmService) ValidateWebhook(ctx context.Context, rawData []byte) error {
	op := "ConfirmService.ValidateWebhook"

	data, err := c.validate.ValidateData(rawData)
	if err != nil {
		c.log.Error("failed to validate data", slog.String("op", op), sl.Err(err))
		return err // temp
	}

	log := c.log.With(
		slog.String("op", op),
		slog.String("user_id", data.UserID),
	)

	log.Info("success validate data!")

	// обращаемся к базе на поиск idempotency_key, если его нету идем дальше
	check, err := c.paymentprv.IdempotencyAndStatus(ctx, data.IdempotencyKey)

	// Оутбокс паттерн создаем, создаем воркер, сначала вносим в 2 таблице,
	// воркер из одной будет в кафку крутить, сообщения не теряются

	return nil // temp
}
