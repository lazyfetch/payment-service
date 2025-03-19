package paymentsrv

import (
	"context"
	"errors"
	"log/slog"
	"payment/internal/domain/models"
	rk "payment/internal/lib/robokassa"
)

const (
	yookassa  = "Yookassa"
	robokassa = "Robokassa"
)

type PaymentService struct {
	log *slog.Logger
}

// New is builder function which return *PaymentService struct (А то оно не видно)
func New(log *slog.Logger) *PaymentService {
	return &PaymentService{}
}

func (p *PaymentService) GetPaymentURL(ctx context.Context, req models.Payment) (string, error) {
	const op = "GetPaymentURL"

	log := p.log.With(
		slog.String("op", op),
		slog.String("user_id", req.UserID),
		slog.String("payment_method", req.PaymentMethod),
	)

	log.Info("attemping to generate url")

	switch req.PaymentMethod {
	case robokassa:
		paymentURL, err := rk.GeneratePaymentURL(req)
		if err != nil {
			return "", err // temp
		}
		return paymentURL, nil // temp
	default:
		return "", errors.New("invalid payment method") // temp
	}
}
