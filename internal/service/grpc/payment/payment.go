package paymentsrv

import (
	"context"
	"errors"
	"log/slog"
	"payment/internal/domain/models"
)

const (
	robokassa = "Robokassa"
)

type PaymentService struct {
	log *slog.Logger
}

func New() *PaymentService {
	return &PaymentService{}
}

func (p *PaymentService) GetPaymentURL(ctx context.Context, req models.Payment) (paymentURL string, err error) {
	const op = "GetPaymentURL"

	log := p.log.With(
		slog.String("op", op),
		slog.String("user_id", req.UserID),
		slog.String("payment_method", req.PaymentMethod),
	)

	switch req.PaymentMethod {
	case robokassa:
		log.Info("attemping to generate url")

		// delaem

	default:
		return "", errors.New("invalid payment method")
	}

	return "", nil // TEMP

}
