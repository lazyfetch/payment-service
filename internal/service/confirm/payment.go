package confirmsrv

import (
	"context"
	"fmt"
	"log/slog"
	"payment/internal/govnokassa"
	"payment/internal/lib/logger/sl"
)

type PaymentUpdater interface {
	OutboxUpdatePayment(ctx context.Context, idemKey, userID string) error
}

type PaymentProvider interface {
	IdempotencyAndStatus(ctx context.Context, idempotencyKey string) bool
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

	check := c.paymentprv.IdempotencyAndStatus(ctx, data.IdempotencyKey)
	if !check {
		return fmt.Errorf("") // тут надо логировать лучше, но мне лень
	}

	// Outbox pattern
	if err = c.paymentupdr.OutboxUpdatePayment(ctx, data.IdempotencyKey, data.UserID); err != nil {
		return err
	}

	return nil // temp
}
