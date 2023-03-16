package main

import (
	"context"
	"log"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"

	"counters/internal/config"
	"counters/internal/http"
	"counters/pkg/counter"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("zap logger is not created: %v", err)
	}
	undo := zap.ReplaceGlobals(l)
	defer undo()

	var cfg config.Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		l.Fatal("config reading error", zap.Error(err))
	}

	cms := counter.NewMemoryStore()
	cm := counter.NewManager(cms)

	http.MustRegisterMetrics(prometheus.DefaultRegisterer)

	s := &nethttp.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: http.NewHandler(l, cm),
	}
	go func() {
		if err := s.ListenAndServe(); err != nil && err != nethttp.ErrServerClosed {
			l.Fatal("HTTP server didn't start", zap.Error(err))
		}
	}()
	l.Info("HTTP server started", zap.String("address", cfg.HTTPServer.Addr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		l.Info("shutting down gracefully")
	case <-ctx.Done():
		l.Info("context has terminated")
	}

	if err := s.Shutdown(ctx); err != nil {
		l.Fatal("HTTP server shutdown failed", zap.Error(err))
	}
	l.Info("HTTP server shut down")
}
