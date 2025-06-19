package app

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/samber/oops"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewDB(c Config) (*bun.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.DB.ConnectTimeout)*time.Minute)
	defer cancel()

	pool, err := pgxpool.New(ctx, c.DB.DSN)
	if err != nil {
		return nil, oops.In("NewDB").Errorf("cannot parse config: %w", err)
	}

	pgdb := stdlib.OpenDBFromPool(pool)
	err = pgdb.PingContext(ctx)
	if err != nil {
		return nil, oops.In("NewDB").Errorf("cannot ping database: %w", err)
	}

	db := bun.NewDB(pgdb, pgdialect.New())

	if c.Server.Debug {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return db, nil
}
