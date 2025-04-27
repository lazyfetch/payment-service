package interceptors

import (
	"context"
	"net"
	"payment/proto/gen/payment"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// В редисе реализуем этот интерфейс
type RateLimiter interface {
	Allow(ctx context.Context, ip string, window time.Duration, maxRequests int, banDurations []time.Duration) (bool, error)
}

type LimiterOpts struct {
	Enabled      bool
	Window       time.Duration
	MaxRequests  int
	BanDurations []time.Duration
}

type ValidateOpts struct {
	MaxNameLength    int
	MaxAmount        int64
	MaxMessageLenght int
	PaymentService   string
}

func LimiterInterceptor(lOpts LimiterOpts, limiter RateLimiter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// temp
		// Плохо, потому что такое обрабатывать на стадии
		// Сборки, а не на стадии рабочей функции, поэтому
		// Оборочивать этот момент, и не юзать такое на проде (пока можно)

		if !lOpts.Enabled {
			return handler(ctx, req)
		}

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
		cond, err := limiter.Allow(ctx, ip, lOpts.Window, lOpts.MaxRequests, lOpts.BanDurations)
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

func ValidationInterceptor(vOpts ValidateOpts) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if _, ok := req.(*payment.GetPaymentUrlRequest); !ok {
			return nil, status.Error(codes.InvalidArgument, "invalid request type")
		}
		if err := ValidateGetPaymentUrl(vOpts, req.(*payment.GetPaymentUrlRequest)); err != nil {
			return nil, err
		}
		return handler(ctx, req)

	}
}

func ValidateGetPaymentUrl(vOpts ValidateOpts, req *payment.GetPaymentUrlRequest) error {

	if req.Name == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}
	if len(req.Name) > vOpts.MaxNameLength {
		return status.Error(codes.InvalidArgument, "name is too long max 40 characters allowed")
	}
	if len(req.Description) > vOpts.MaxMessageLenght {
		return status.Error(codes.InvalidArgument, "description is too long, max 250 characters allowed")
	}
	if req.Amount <= 0 {
		if req.Amount >= vOpts.MaxAmount {
			return status.Error(codes.InvalidArgument, "amount exceeds the limit: 1000000")
		}
		return status.Error(codes.InvalidArgument, "amount must be positive")
	}
	if req.PaymentMethod != vOpts.PaymentService {
		return status.Error(codes.InvalidArgument, "no such payment method exists")
	}
	if req.UserId == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}
