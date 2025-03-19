package confirmsrv

import (
	"context"
	"payment/internal/domain/models"
)

type PaymentSaver interface {
	CreatePayment(ctx context.Context, data models.Payment)
}

type ConfirmService struct{}
