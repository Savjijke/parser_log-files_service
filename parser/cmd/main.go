package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/savjijke/parser-log-files-service/internal/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()
	cfg := config.MustLoad(configPath)
	log := mustMakeLogger(cfg.LogLevel)
	fmt.Println(log) // UBRAT' POTOM
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		slog.Error("error unmarshaling loglevel", "msg error: ", err)
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}
