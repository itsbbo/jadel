package model

import "github.com/jackc/pgx/v5/pgconn"

func IsErrUniqueEmailUser(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}

	return pgErr.Code == "23505" && pgErr.TableName == "users" && pgErr.ColumnName == "email"
}

func IsErrConstrainPrivateKeyInServer(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}

	return pgErr.Code == "23503" && pgErr.TableName == "servers" && pgErr.ColumnName == "private_key_id"
}
