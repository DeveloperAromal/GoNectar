package httptrap

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/DeveloperAromal/GoNectar/internal/collector"
	"github.com/DeveloperAromal/GoNectar/internal/config"
)

type HTTPTrap struct {
	srv *http.Server
	log *log.Logger
}

func NewHTTPTrap(cfg *config.Config, col *collector.Collector, logger *log.Logger) *HTTPTrap {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		_ = r.Body.Close()

		remote := r.RemoteAddr

		if x := r.Header.Get("X-Forwarded-For"); x != "" {
			remote = strings.Split(x, ",")[0]
		}

		col.IngestEvent(collector.Event{
			Type: "http.request",
			Time: time.Now().UTC(),
			Date: map[string]interface{}{
				"method":  r.Method,
				"path":    r.URL.Path,
				"remote":  remote,
				"ua":      r.UserAgent(),
				"headers": r.Header,
				"body":    string(body),
			},
		})

		if r.URL.Path == "/login" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`
									<!DOCTYPE html>
									<html lang="en" class="bg-black text-white">
									<head>
										<meta charset="UTF-8" />
										<meta name="viewport" content="width=device-width, initial-scale=1.0" />
										<title>Document</title>
										<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
									</head>
									<body class="flex items-center justify-center h-screen">
										<div>
										<div><h2 class="text-amber-400 text-3xl">你被抓住了</h2></div>
										</div>
									</body>
									</html>
									`))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Welcome"))
	})

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &HTTPTrap{srv: srv, log: logger}
}

func (h *HTTPTrap) Start() error {
	ln, err := net.Listen("tcp", h.srv.Addr)
	if err != nil {
		return err
	}
	h.log.Printf("HTTP trap starting on %s\n", h.srv.Addr)
	if err := h.srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (h *HTTPTrap) Stop(ctx context.Context) error {
	h.log.Println("stopping HTTP trap")
	if err := h.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
