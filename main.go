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
			slog.Error(msg, r.Host, r.URL.Path, r.Method, r.RemoteAddr)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		b := string(d)

		slog.Info(fmt.Sprintf("request received:\n%s\n\n", b), r.Host, r.URL.Path, r.Method, r.RemoteAddr)
		if _, err := fmt.Fprint(w, b); err != nil {
			msg := fmt.Sprintf("couldn't write response: %s", err)
			slog.Error(msg, r.Host, r.URL.Path, r.Method, r.RemoteAddr)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	})

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	addr := ":" + port

	slog.Info(fmt.Sprintf("http-dump is listening at %s\n", addr))
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error(fmt.Sprintf("couldn't start server: %s", err))
	}
	slog.Info("http-dump is shutting down")
}
