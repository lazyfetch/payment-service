package confirmsrv

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"payment/internal/govnokassa"
	"payment/internal/lib/logger/sl"
)

type PaymentUpdater interface {
	OutboxUpdatePayment(ctx context.Context, idemKey string, payload any) error
}

type PaymentProvider interface {
	IdempotencyAndStatus(ctx context.Context, idempotencyKey string) error
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
	const op = "ConfirmService.ValidateWebhook"

	data, err := c.validate.ValidateData(rawData)
	if err != nil {
		c.log.Error("failed to validate data", slog.String("op", op), sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log := c.log.With(
		slog.String("op", op),
		slog.String("user_id", data.UserID),
	)

	log.Info("success validate data!")

	if err = c.paymentprv.IdempotencyAndStatus(ctx, data.IdempotencyKey); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("in_progress and idempotency key not found")
			return fmt.Errorf("in_progress and idempotency key not found")
		}
		log.Error("failed to check idem_key and status")
		return fmt.Errorf("%s: %w", op, err)
	}

	// Outbox pattern
	if err = c.paymentupdr.OutboxUpdatePayment(ctx, data.IdempotencyKey, data); err != nil {
		log.Error("failed to update payment", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("success webhook validation!")

	return nil
}
