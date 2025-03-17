package paymentsrv

import (
	"context"
	"payment/internal/domain/models"
)

type PaymentService struct {
}

func New() *PaymentService {
	return &PaymentService{}
}

func (p *PaymentService) GetPaymentUrl(context.Context, models.PaymentRequest) string {

	return "" // temp
}
