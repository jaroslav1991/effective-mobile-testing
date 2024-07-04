package connection

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
)

func NewPostgresDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		slog.Error("can't connect to DB:", slog.String("err", err.Error()))
		return nil, err
	}
	db.MustBegin()

	return db, nil
}

func InitSchema(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		slog.Error("can't initialize postgres driver:", slog.String("err", err.Error()))
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file:///go/src/app/migrations", "postgres", driver)
	if err != nil {
		slog.Error("can't migrate:", slog.String("err", err.Error()))
		return err
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("no change to migration history")
			return nil
		}
		slog.Error("can't up migrate:", slog.String("err", err.Error()))
		return err
	}

	return nil
}
