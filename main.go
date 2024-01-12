package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

var port = "8080"

type structedRequest struct {
	Proto  string
	Method string
	Host   string
	Header http.Header
}

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: time.RFC3339,
	}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var s structedRequest
		s.Proto = r.Proto
		s.Method = r.Method
		s.Host = r.Host
		s.Header = r.Header

		d, err := httputil.DumpRequest(r, true)
		if err != nil {
			msg := fmt.Sprintf("couldn't dump request: %v", err)
			logger.Error(msg, slog.Any("request", s))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		b := string(d)
		logger.Info("request received", slog.Any("request", s))

		if _, err := fmt.Fprint(w, b); err != nil {
			msg := fmt.Sprintf("couldn't write response: %s", err)
			logger.Error(msg, slog.Any("request", s))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	})

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	addr := ":" + port
	logger.Info("http-dump is starting", "addr", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Error(fmt.Sprintf("couldn't start server: %s", err))
	}
	logger.Info("http-dump is shutting down")
}
