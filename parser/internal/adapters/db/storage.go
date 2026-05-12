package db

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

const locationForLogger = "adapters/db/"

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

func NewDB(dsn string, log *slog.Logger) (*DB, error) {

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Error("connection problem", "address", dsn, "error", err)
		return nil, err
	}
	log.Debug("DB connected", "db dsn", dsn)
	return &DB{log: log, conn: db}, nil
}
