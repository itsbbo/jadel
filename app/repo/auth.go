package repo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"time"

	"github.com/itsbbo/jadel/app/auth"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/drivers/pgx"
)

type Auth struct {
	db pgx.Pool
}

func NewAuth(db pgx.Pool) *Auth {
	return &Auth{db: db}
}

func (d *Auth) NewUserWithSession(ctx context.Context, r auth.NewUserWithSessionParam) (*model.User, string, error) {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return nil, "", err
	}

	defer tx.Rollback(ctx)

	userId := ulid.Make()
	userSetter := model.UserSetter{
		ID:       &userId,
		Name:     &r.Name,
		Email:    &r.Email,
		Password: &r.Password,
	}

	user, err := model.Users.Insert(&userSetter).One(ctx, tx)
	if err != nil {
		return nil, "", err
	}

	sessionID := rand.Text()
	sessionsExpired := time.Now().Add(3 * time.Hour)

	user.InsertSessions(ctx, d.db, &model.SessionSetter{
		ID:        &sessionID,
		UserID:    &userId,
		ExpiredAt: &sessionsExpired,
	})

	err = tx.Commit(ctx)
	if err != nil {
		return nil, "", err
	}

	return user, sessionID, nil
}

func (d *Auth) FindByEmailPassword(ctx context.Context, email string, password string) (*model.User, error) {
	user, err := model.Users.Query(
		psql.WhereAnd(
			model.SelectWhere.Users.Email.EQ(email),
			model.SelectWhere.Users.Password.EQ(psql.Raw(`crypt($2, password)`, password).String()),
		),
	).One(ctx, d.db)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *Auth) InsertSession(ctx context.Context, param auth.InsertSessionParam) error {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	expiredAt := time.Now().Add(param.Expires)
	ipaddr := sql.Null[string]{
		V: param.IPAddr,
		Valid: len(param.IPAddr) != 0,
	}

	userAgent := sql.Null[string]{
		V: param.UserAgent,
		Valid: len(param.UserAgent) != 0,
	}

	_, err = model.Sessions.Insert(&model.SessionSetter{
		ID:        &param.SessionID,
		UserID:    &param.UserID,
		ExpiredAt: &expiredAt,
		IPAddress: &ipaddr,
		UserAgent: &userAgent,
	}).Exec(ctx, tx)

	return tx.Commit(ctx)
}
