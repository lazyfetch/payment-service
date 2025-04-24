package interceptors

import (
	"context"
	"net"
	"payment/internal/config"
	"payment/proto/gen/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// В редисе реализуем этот интерфейс
type RateLimiter interface {
	Allow(ctx context.Context, ip string) (bool, error)
}

func LimiterInterceptor(cfg *config.Internal, limiter RateLimiter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// Получаем socket
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "something wrong") // temp
		}
		// Парсим ip
		addr := p.Addr.String()
		ip, err := extractIP(addr)
		if err != nil {
			return nil, status.Error(codes.Internal, "something wrong")
		}
		// Получаем условие
		cond, err := limiter.Allow(ctx, ip)
		if err != nil {
			return nil, status.Error(codes.Internal, "something wrong") // temp
		}
		// Проверяем условие
		if !cond {
			return nil, status.Error(codes.ResourceExhausted, "too many requests")
		}
		// Если все ок возвращаем хендлер, и пропускаем этого мамкиного спамера
		return handler(ctx, req)
	}
}

func extractIP(addr string) (string, error) {
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	return ip, nil
}

func ValidationInterceptor(cfg *config.Internal) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if _, ok := req.(*payment.GetPaymentUrlRequest); !ok {
			return nil, status.Error(codes.InvalidArgument, "invalid request type")
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
