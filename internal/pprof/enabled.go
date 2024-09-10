//go:build pprof

package pprof

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"os"
	"time"
)

func ListenAndServe() {
	go func() {
		addr := "127.0.0.1:6060"
		if env := os.Getenv("PPROF_ADDR"); env != "" {
			addr = env
		}

		server := &http.Server{
			Addr:        addr,
			ReadTimeout: 10 * time.Second,
		}
		slog.Info("Starting pprof server", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			slog.Error("Failed to start pprof server", "error", err.Error())
		}
	}()
}
