package paymentgrpc

import (
	"context"
	"errors"
	"payment/internal/domain/models"
	generatesrv "payment/internal/service/generate"
	payment "payment/proto/gen/payment"

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

	// validate
	if err := validateGetPaymentUrl(req); err != nil {
		return nil, err
	}

	// generate
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

func validateGetPaymentUrl(req *payment.GetPaymentUrlRequest) error {

	if len(req.Name) > 40 {
		if req.Name == "" {
			return status.Error(codes.InvalidArgument, "name is required")
		}
		return status.Error(codes.InvalidArgument, "name is too long max 40 characters allowed") // костыль, в конфиг надо temp
	}
	if len(req.Description) > 250 {
		return status.Error(codes.InvalidArgument, "description is too long, max 250 characters allowed") // костыль, в конфиг надо temp
	}
	if req.Amount <= 0 {
		if req.Amount >= 100000000 { // желательно через конфиг передавать, temp
			return status.Error(codes.InvalidArgument, "amount exceeds the limit: 1000000")
		}
		return status.Error(codes.InvalidArgument, "amount must be positive")
	}
	if req.PaymentMethod != "Robokassa" {
		return status.Error(codes.InvalidArgument, "no such payment method exists")
	}
	if req.UserId == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}
