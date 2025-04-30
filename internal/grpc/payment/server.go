package paymentgrpc

import (
	"context"
	"errors"
	"payment/internal/domain/models"
	generatesrv "payment/internal/service/generate"
	payment "payment/proto/gen/payment"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	tracer := otel.Tracer("payment-service")

	ctx, span := tracer.Start(ctx, "create-payment")
	defer span.End()

	url, err := s.payment.GetPaymentURL(ctx, models.GRPCPayment{
		Name:          req.Name,
		Description:   req.Description,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		UserID:        req.UserId,
	})
	if err != nil {
		if errors.Is(err, generatesrv.ErrInvalidUserID) {
			return nil, status.Error(codes.InvalidArgument, "user_id not found")
		}
		if errors.Is(err, generatesrv.ErrAmountTooSmall) {
			return nil, status.Error(codes.InvalidArgument, "amount is too smail")
		}
		return nil, status.Error(codes.Internal, "failed to generate url")
	}

	return &payment.GetPaymentUrlResponse{
		PaymentUrl: url,
	}, nil

}
