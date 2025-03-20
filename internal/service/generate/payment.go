package generatesrv

import (
	"context"
	"errors"
	"log/slog"
	"payment/internal/domain/models"
)

var (
	ErrInvalidArgument = errors.New("invalid amount")
)

type GeneratePaymentURL interface {
	GeneratePaymentURL(models.GRPCPayment) (string, error)
}

type PaymentSaver interface {
	CreatePayment(ctx context.Context, data models.DBPayment) error
}

type PaymentService struct {
	log        *slog.Logger
	paymentsvr PaymentSaver
	paymentgen GeneratePaymentURL
}

// New is builder function which return *PaymentService struct (А то оно не видно)
func New(log *slog.Logger, paymentsvr PaymentSaver, paymentgen GeneratePaymentURL) *PaymentService {
	return &PaymentService{
		log:        log,
		paymentsvr: paymentsvr,
		paymentgen: paymentgen,
	}
}

func (p *PaymentService) GetPaymentURL(ctx context.Context, req models.GRPCPayment) (string, error) {
	const op = "paymentService.GetPaymentURL"

	log := p.log.With(
		slog.String("op", op),
		slog.String("user_id", req.UserID),
		slog.String("payment_method", req.PaymentMethod),
	)

	log.Info("Attemping to generate url")

	// Чекаем наличие такого юзера, возвращаем ошибку если такого нету

	// Чекаем его минималку донатную, возвращаем ошибку если ниже минималки

	paymentURL, err := p.paymentgen.GeneratePaymentURL(req)
	if err != nil {
		return "", err // temp
	}

	return paymentURL, nil
}
