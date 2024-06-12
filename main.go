package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: time.RFC3339,
	}))
	slog.SetDefault(logger)

	ctx := context.Background()

	close, err := NewTracerProvider(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create tracer provider: %v", err))
	}
	defer close()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path))

		s := NewStructedRequest(r)

		d, err := httputil.DumpRequest(r, true)
		if err != nil {
			msg := fmt.Sprintf("couldn't dump request: %v", err)
			span.RecordError(err)
			slog.Error(msg, slog.Any("request", s))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		span.AddEvent("request received", trace.WithAttributes(
			s.GetAttributes()...,
		))

		slog.Info("request received", slog.Any("request", s))

		if _, err := fmt.Fprint(w, string(d)); err != nil {
			msg := fmt.Sprintf("couldn't write response: %s", err)
			span.RecordError(err)
			slog.Error(msg, slog.Any("request", s))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		defer span.End()
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	slog.Info("http-dump is starting", "addr", addr)

	if err := http.ListenAndServe(addr, otelhttp.NewHandler(mux, "server",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)); err != nil {
		slog.Error(fmt.Sprintf("couldn't start server: %s", err))
	}
	slog.Info("http-dump is shutting down")
}
