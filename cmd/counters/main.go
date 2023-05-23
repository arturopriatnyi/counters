package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"counters/internal/config"
	"counters/internal/handler"
	"counters/pkg/counter"
	"counters/pkg/iam"
	"counters/pkg/logger"
	"counters/pkg/oauth2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cfg config.Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Fatalf("config reading failed: %v", err)
	}

	l, err := logger.New(logger.WithMode(logger.Mode(cfg.Mode)))
	if err != nil {
		log.Fatalf("zap logger creating failed: %v", err)
	}
	undo := zap.ReplaceGlobals(l)
	defer undo()

	l.Info("starting", zap.Any("mode", cfg.Mode))

	cms := counter.NewMemoryStorage()
	cm := counter.NewManager(cms)

	ums := iam.NewUserMemoryStorage()
	iamm := iam.NewManager(ums, map[oauth2.Provider]oauth2.Client{
		oauth2.Google: oauth2.NewGoogleClient(
			oauth2.Config{
				ClientID:     cfg.GoogleOAuth2.ClientID,
				ClientSecret: cfg.GoogleOAuth2.ClientSecret,
				RedirectURL:  cfg.GoogleOAuth2.RedirectURL,
				Scopes:       cfg.GoogleOAuth2.Scopes,
			},
			http.DefaultClient,
		),
		oauth2.GitHub: oauth2.NewGitHubClient(
			oauth2.Config{
				ClientID:     cfg.GitHubOAuth2.ClientID,
				ClientSecret: cfg.GitHubOAuth2.ClientSecret,
				RedirectURL:  cfg.GitHubOAuth2.RedirectURL,
				Scopes:       cfg.GitHubOAuth2.Scopes,
			},
			http.DefaultClient,
		),
	})
	handler.MustRegisterMetrics(prometheus.DefaultRegisterer)

	s := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: handler.New(l, iamm, cm),
	}
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal("HTTP server starting failed", zap.Error(err))
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
