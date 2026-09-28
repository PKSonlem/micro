package logger

import (
	"context"
	"log/slog"
	"os"
)

var (
	defaultTagName tagName = "default"
)

type tagName string

type field struct {
	f map[string]string
}

type contextHandler struct {
	slog.Handler
}

func newContextHandler(h slog.Handler) contextHandler {
	return contextHandler{h}
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	t, _ := ctx.Value(defaultTagName).(field)

	attrs := make([]slog.Attr, len(t.f))
	for idx, value := range t.f {
		attrs = append(attrs, slog.String(idx, value))
	}
	r.AddAttrs(attrs...)

	return h.Handler.Handle(ctx, r)
}

type Logger struct {
	*slog.Logger
}

func New() *Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	ctxHandler := newContextHandler(jsonHandler)
	return &Logger{slog.New(&ctxHandler)}
}

func buildFields() field {
	m := make(map[string]string)
	return field{f: m}
}

func (l *Logger) WithFields(ctx context.Context, f map[string]string) context.Context {
	var t field

	t, ok := ctx.Value(defaultTagName).(field)
	if !ok {
		t = buildFields()
	}

	for field, val := range f {
		t.f[field] = val
	}

	return context.WithValue(ctx, defaultTagName, t)
}

func (l *Logger) Info(ctx context.Context, message string, args ...any) {
	l.Info(ctx, message, args...)
}

func (l *Logger) Error(ctx context.Context, err error, args ...any) {
	l.Error(ctx, err, args...)
}
