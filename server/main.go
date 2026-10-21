// Command portfolio — HTTP-сервер портфолио: JSON API с контентом + раздача
// собранного Vue-фронтенда из встроенной файловой системы.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// version подставляется линкером:
//
//	go build -ldflags "-X main.version=$(git rev-parse --short HEAD)"
var version = "dev"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg := loadConfig()

	store, err := NewStore(cfg.DataDir, newRemoteStore())
	if err != nil {
		logger.Error("не удалось открыть хранилище контента", "err", err)
		os.Exit(1)
	}

	auth, err := NewAuth(cfg.DataDir)
	if err != nil {
		logger.Error("не удалось настроить вход в админку", "err", err)
		os.Exit(1)
	}

	app := newApp(version, store, auth)

	mux := http.NewServeMux()
	app.routes(mux)

	handler := chain(mux,
		recoverPanic,
		requestLogger,
		securityHeaders,
		cors(cfg.DevOrigin),
	)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)
		<-ctx.Done()
		logger.Info("остановка сервера", "timeout", cfg.ShutdownTimeout.String())
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("graceful shutdown не удался", "err", err)
		}
	}()

	logger.Info("портфолио запущено", "addr", "http://localhost:"+cfg.Port, "version", version, "dev_origin", cfg.DevOrigin)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("сервер упал", "err", err)
		os.Exit(1)
	}

	<-shutdownDone
	logger.Info("сервер остановлен")
}

type config struct {
	Port            string
	DataDir         string
	DevOrigin       string // пусто = CORS выключен (прод: один origin)
	ShutdownTimeout time.Duration
}

func loadConfig() config {
	return config{
		Port:            envOr("PORT", "8080"),
		DataDir:         envOr("DATA_DIR", "data"),
		DevOrigin:       os.Getenv("DEV_ORIGIN"),
		ShutdownTimeout: 10 * time.Second,
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
