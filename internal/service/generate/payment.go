package generatesrv

import (
	"context"
	"errors"
	"log/slog"
	"payment/internal/domain/models"
	"payment/internal/lib/logger/sl"
	"payment/internal/storage"
)

var (
	ErrInvalidUserID  = errors.New("invalid user_id")
	ErrAmountTooSmall = errors.New("amount is too small")
)

type UserProvider interface {
	GetMinAmountByUser(ctx context.Context, userID string) (int64, error)
}

type GeneratePaymentURL interface {
	GeneratePaymentURL(models.GRPCPayment) (string, error)
}

type PaymentSaver interface {
	CreatePayment(ctx context.Context, data models.DBPayment) error
}

type PaymentService struct {
	log        *slog.Logger
	paymentsvr PaymentSaver
	userprv    UserProvider
	paymentgen GeneratePaymentURL
}

// New is builder function which return *PaymentService struct (А то оно не видно)
func New(log *slog.Logger, paymentsvr PaymentSaver, userprv UserProvider, paymentgen GeneratePaymentURL) *PaymentService {
	return &PaymentService{
		log:        log,
		paymentsvr: paymentsvr,
		userprv:    userprv,
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

	log.Info("attemping to generate url")

	minAmount, err := p.userprv.GetMinAmountByUser(ctx, req.UserID)

	if err != nil {
		if errors.Is(err, storage.ErrUserIDNotFound) {
			log.Warn("user_id not found")
			return "", ErrInvalidUserID
		}
		log.Error("failed to check min_amount", sl.Err(err))
	}

	if req.Amount < minAmount {
		log.Warn("min_amount too small")
		return "", ErrAmountTooSmall
	}

	// создаем UUID имплементим

	// передаем в GOVNOKASSA mock edition генератор

	paymentURL, err := p.paymentgen.GeneratePaymentURL(req)
	if err != nil {
		return "", err // temp
	}

	// Если нету ошибок, записываем в бд

	return paymentURL, nil
}
