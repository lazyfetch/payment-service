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
		return fmt.Errorf("failed to validate webhook")
	}

	log := c.log.With(
		slog.String("op", op),
		slog.String("user_id", data.UserID),
	)

	log.Info("success validate data!")

	check := c.paymentprv.IdempotencyAndStatus(ctx, data.IdempotencyKey)
	if !check {
		log.Error("failed to check idem_key and status")
		return fmt.Errorf("failed to check idem_key and status")
	}

	// Outbox pattern
	if err = c.paymentupdr.OutboxUpdatePayment(ctx, data.IdempotencyKey, data.UserID); err != nil {
		log.Error("failed to update payment", sl.Err(err))
		return fmt.Errorf("failed to update payment")
	}

	log.Info("success webhook validation!")

	return nil
}
