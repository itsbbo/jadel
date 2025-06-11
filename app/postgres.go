package app

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob/drivers/pgx"
)

func NewPostgres(c Config) (pgx.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.DB.ConnectTimeout))
	defer cancel()

	db, err := pgx.New(ctx, c.DB.DSN)
	if err != nil {
		return pgx.Pool{}, err
	}

	err = db.Ping(ctx)

	return db, err
}
