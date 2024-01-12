package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
)

var port = "8080"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		d, err := httputil.DumpRequest(r, true)
		if err != nil {
			msg := fmt.Sprintf("couldn't dump request: %v", err)
			slog.Error(msg, "host", slog.Any("request", r))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		b := string(d)
		slog.Info("request received", slog.Any("request", r))

		if _, err := fmt.Fprint(w, b); err != nil {
			msg := fmt.Sprintf("couldn't write response: %s", err)
			slog.Error(msg, slog.Any("request", r))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	})

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	addr := ":" + port
	slog.Info("http-dump is starting", "addr", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error(fmt.Sprintf("couldn't start server: %s", err))
	}
	slog.Info("http-dump is shutting down")
}
