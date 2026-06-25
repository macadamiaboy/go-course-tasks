package observability

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type TraceHandler struct {
	slog.Handler
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx == nil {
		return h.Handler.Handle(ctx, r)
	}

	span := trace.SpanFromContext(ctx)
	if spanCtx := span.SpanContext(); spanCtx.IsValid() {
		r.AddAttrs(
			slog.String("trace_id", spanCtx.TraceID().String()),
			slog.String("span_id", spanCtx.SpanID().String()),
		)
	}

	return h.Handler.Handle(ctx, r)
}

func (h *TraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TraceHandler{
		Handler: h.Handler.WithAttrs(attrs),
	}
}

func (h *TraceHandler) WithGroup(name string) slog.Handler {
	return &TraceHandler{
		Handler: h.Handler.WithGroup(name),
	}
}

func NewLogger(service, env, version string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(&TraceHandler{handler}).With("service", service, "env", env, "version", version)
}
