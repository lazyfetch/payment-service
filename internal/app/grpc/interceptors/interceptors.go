package interceptors

import (
	"context"
	"payment/internal/config"
	"payment/proto/gen/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidationInterceptor(cfg *config.Internal) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if _, ok := req.(*payment.GetPaymentUrlRequest); !ok {
			return nil, status.Error(codes.InvalidArgument, "Invalid request type")
		}
		if err := ValidateGetPaymentUrl(cfg, req.(*payment.GetPaymentUrlRequest)); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func ValidateGetPaymentUrl(cfg *config.Internal, req *payment.GetPaymentUrlRequest) error {

	if req.Name == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}
	if len(req.Name) > cfg.MaxNameLength {
		return status.Error(codes.InvalidArgument, "name is too long max 40 characters allowed")
	}
	if len(req.Description) > cfg.MaxMessageLenght {
		return status.Error(codes.InvalidArgument, "description is too long, max 250 characters allowed")
	}
	if req.Amount <= 0 {
		if req.Amount >= cfg.MaxAmount {
			return status.Error(codes.InvalidArgument, "amount exceeds the limit: 1000000")
		}
		return status.Error(codes.InvalidArgument, "amount must be positive")
	}
	if req.PaymentMethod != cfg.PaymentService {
		return status.Error(codes.InvalidArgument, "no such payment method exists")
	}
	if req.UserId == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}
