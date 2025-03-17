package paymentgrpc

import (
	"context"
	"payment/internal/domain/models"
	payment "payment/proto/payment/gen"

	"google.golang.org/grpc"
)

type PaymentService interface {
	GetPaymentUrl(ctx context.Context, req models.PaymentRequest) string
}

type serverAPI struct {
	payment.UnimplementedPaymentServer
	payment PaymentService
}

func Register(gRPC *grpc.Server, paymentService PaymentService) {
	payment.RegisterPaymentServer(gRPC, &serverAPI{payment: paymentService})
}

func (s *serverAPI) GetPaymentUrl(ctx context.Context, req *payment.GetPaymentUrlRequest) (*payment.GetPaymentUrlResponse, error) {

	s.payment.GetPaymentUrl(ctx, models.PaymentRequest{})

	return &payment.GetPaymentUrlResponse{
		PaymentLink: "",
	}, nil
}
