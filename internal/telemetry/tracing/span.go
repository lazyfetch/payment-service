package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Span interface {
	End()
	SetAttributes(attrs ...attribute.KeyValue)
	RecordError(err error)
	Context() context.Context
	Name() string
}

type spanImpl struct {
	ctx  context.Context
	span trace.Span
	name string
}

func StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, Span) {
	tracer := otel.Tracer("payment-service") // тут можно сделать динамическим
	ctx, span := tracer.Start(ctx, name)
	span.SetAttributes(attrs...)

	return ctx, &spanImpl{
		ctx:  ctx,
		span: span,
		name: name,
	}
}

func (s *spanImpl) End() {
	s.span.End()
}

func (s *spanImpl) SetAttributes(attrs ...attribute.KeyValue) {
	s.span.SetAttributes(attrs...)
}

func (s *spanImpl) RecordError(err error) {
	if err != nil {
		s.span.RecordError(err)
		s.span.SetStatus(codes.Error, err.Error())
	}
}

func (s *spanImpl) Context() context.Context {
	return s.ctx
}

func (s *spanImpl) Name() string {
	return s.name
}
