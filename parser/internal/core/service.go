package core

import (
	"archive/zip"
	"context"
	"log/slog"
	"path/filepath"
	"sync/atomic"
)

const locationForLogger = "core/"

type Service struct {
	db      DB
	log     *slog.Logger
}

func NewService(db DB, log *slog.Logger) *Service {
	return &Service{db: db, log: log}
}

// парсить параллельно .db_csv по start, проверять что кол-во старт и енд равно и больше 0
func (s *Service) Parse(ctx context.Context, url string) (int, error) {

}

