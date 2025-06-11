package app

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/samber/oops"
	"github.com/stephenafamo/bob/drivers/pgx"
)

func NewPostgres(c Config) (pgx.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.DB.ConnectTimeout)*time.Minute)
	defer cancel()

	db, err := pgx.New(ctx, c.DB.DSN)
	if err != nil {
		return pgx.Pool{}, oops.In("NewPostgres").Errorf("cannot create new postgres: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return db, oops.In("NewPostgres").Errorf("cannot ping database: %w", err)
	}

	return db, nil
}
