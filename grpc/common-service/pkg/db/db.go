package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.nhat.io/otelsql"
)

func InitDB(ctx context.Context, driver, dsn string) (*sql.DB, error) {
	driverName, err := otelsql.Register("postgres",
		otelsql.AllowRoot(),
		otelsql.TraceQueryWithoutArgs(),
		otelsql.TraceRowsClose(),
		otelsql.TraceRowsAffected(),
	)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	// Verify DB is up
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func ApplyMigrations(db *sql.DB, driver, migrationsDir string) error {
	if err := goose.SetDialect(driver); err != nil {
		return err
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}
	return nil
}
