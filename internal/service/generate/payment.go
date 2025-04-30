package generatesrv

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"payment/internal/domain/models"
	"payment/internal/lib/logger/sl"
	"payment/internal/lib/uuid"
	"payment/internal/storage"
	"payment/internal/telemetry/tracing"

	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrInvalidUserID  = errors.New("invalid user_id")
	ErrAmountTooSmall = errors.New("amount is too small")
)

type UserProvider interface {
	GetMinAmountWithCache(ctx context.Context, userID string) (int64, error)
}

type GeneratePaymentURL interface {
	GeneratePaymentURL(ctx context.Context, data *models.DBPayment) (string, error)
}

type PaymentSaver interface {
	CreatePayment(ctx context.Context, data *models.DBPayment) error
}

type PaymentService struct {
	log        *slog.Logger
	paymentsvr PaymentSaver
	userprv    UserProvider
	paymentgen GeneratePaymentURL
}

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

	ctx, span := tracing.StartSpan(ctx, "Service GetPaymentUrl",
		attribute.String("user_id", req.UserID),
		attribute.String("payment_method", req.PaymentMethod))

	defer span.End()

	log := p.log.With(
		slog.String("op", op),
		slog.String("user_id", req.UserID),
		slog.String("payment_method", req.PaymentMethod),
	)

	log.Info("attemping to generate url")

	minAmount, err := p.userprv.GetMinAmountWithCache(ctx, req.UserID)

	if err != nil {
		if errors.Is(err, storage.ErrUserIDNotExists) {
			span.RecordError(err)
			log.Error("user_id not found")
			return "", ErrInvalidUserID
		}
		span.RecordError(err)
		log.Error("failed to check min_amount", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if req.Amount < minAmount {
		span.RecordError(ErrAmountTooSmall)
		log.Error("min_amount too small")
		return "", ErrAmountTooSmall
	}

	// маппуем, создаем idempotency_key
	idempotencyKey := uuid.UUID()
	payment := models.MapGRPCToDB(&req, idempotencyKey)

	// передаем в GOVNOKASSA edition генератор
	paymentURL, err := p.paymentgen.GeneratePaymentURL(ctx, payment)
	if err != nil {
		span.RecordError(err)
		log.Error("failed to create payment url", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// записываем в бд наш созданный платеж
	if err := p.paymentsvr.CreatePayment(ctx, payment); err != nil {
		span.RecordError(err)
		log.Error("failed to create payment", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// если ошибок нету, вернется ссылка, и кайфарик будет плотный
	log.Info("success!", slog.String("idempotency_key", idempotencyKey))
	return paymentURL, nil
}
