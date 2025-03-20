package paymentgrpc

import (
	"context"
	"fmt"
	"payment/internal/domain/models"
	payment "payment/proto/gen/payment"

	"google.golang.org/grpc"
)

type PaymentService interface {
	GetPaymentURL(ctx context.Context, req models.GRPCPayment) (paymentURL string, err error)
}

type serverAPI struct {
	payment.UnimplementedPaymentServer
	payment PaymentService
}

func Register(gRPC *grpc.Server, paymentService PaymentService) {
	payment.RegisterPaymentServer(gRPC, &serverAPI{payment: paymentService})
}

func (s *serverAPI) GetPaymentUrl(ctx context.Context, req *payment.GetPaymentUrlRequest) (*payment.GetPaymentUrlResponse, error) {

	fmt.Println("HEI")
	url, err := s.payment.GetPaymentURL(ctx, models.GRPCPayment{
		Name:          req.GetName(),
		Description:   req.GetDescription(),
		Amount:        req.GetAmount(),
		PaymentMethod: req.GetPaymentMethod(),
		UserID:        req.GetUserId(),
	})

	if err != nil {
		return nil, err // temp
	}

	return &payment.GetPaymentUrlResponse{
		PaymentUrl: url,
	}, nil
}
