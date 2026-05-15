package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	csvparser "github.com/savjijke/parser-log-files-service/internal/adapters/csvParser"
	"github.com/savjijke/parser-log-files-service/internal/adapters/db"
	"github.com/savjijke/parser-log-files-service/internal/adapters/rest"
	"github.com/savjijke/parser-log-files-service/internal/config"
	"github.com/savjijke/parser-log-files-service/internal/core"
)

func main() {
	cfg := config.MustLoad()
	log := mustMakeLogger(cfg.LogLevel)
	if err := run(cfg, log); err != nil {
		log.Error("server failed", "err", err)
		os.Exit(1)
	}

}

func run(cfg config.Config, log *slog.Logger) error {
	log.Info("starting server")
	log.Debug("debug messages are enabled")

	parser := csvparser.NewParser(log)
	storage, err := db.NewDB(cfg.DataBaseUrl, log)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %v", err)
	}
	if err := storage.Migrate(); err != nil {
		return fmt.Errorf("failed to migrate db: %v", err)
	}
	parserService := core.NewService(storage, parser, log)

	mux := http.NewServeMux()

	mux.Handle("POST /api/v1/parse/", rest.NewParseHandler(log, parserService))
	mux.Handle("GET /api/v1/topology/", rest.NewTopologyHandler(log, parserService))
	mux.Handle("GET /api/v1/node/", rest.NewNodeHandler(log, parserService))
	mux.Handle("GET /api/v1/port/", rest.NewPortHandler(log, parserService))
	mux.Handle("GET /api/v1/log/", rest.NewLogHandler(log, parserService))

	addr := fmt.Sprintf(":%d", cfg.Port)

	server := http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Debug("shutting down server")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error("erroneous shutdown", "error", err)
		}
	}()

	log.Info("Running HTTP server", "address", addr)
	err = server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Error("server closed unexpectedly", "error", err)
		return err
	}
	return nil
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		slog.Error("error unmarshaling loglevel", "err", err)
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}
