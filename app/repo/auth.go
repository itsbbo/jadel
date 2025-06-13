package repo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/app/auth"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
	"github.com/samber/oops"
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

	userSetter := model.UserSetter{
		ID:       omit.From(ulid.Make()),
		Name:     omit.From(r.Name),
		Email:    omit.From(r.Email),
		Password: omit.From(r.Password),
	}

	user, err := model.Users.Insert(&userSetter).One(ctx, tx)
	if err != nil {
		return nil, "", oops.In("model.Users.Insert").Errorf("failed to insert user: %w", err)
	}

	session, err := model.Sessions.Insert(&model.SessionSetter{
		ID:        omit.From(rand.Text()),
		UserID:    omit.From(user.ID),
		IPAddress: omitnull.From(r.IPAddr),
		UserAgent: omitnull.From(r.UserAgent),
		ExpiredAt: omit.From(time.Now().Add(3 * time.Hour)),
	}).One(ctx, tx)

	if err != nil {
		return nil, "", oops.In("user.InsertSessions").Errorf("failed to insert user: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, "", oops.In("tx.Commit").Errorf("failed to commit transaction: %w", err)
	}

	return user, session.ID, nil
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

	_, err = model.Sessions.Insert(&model.SessionSetter{
		ID:        omit.From(param.SessionID),
		UserID:    omit.From(param.UserID),
		ExpiredAt: omit.From(expiredAt),
		IPAddress: omitnull.From(param.IPAddr),
		UserAgent: omitnull.From(param.UserAgent),
	}).Exec(ctx, tx)

	return tx.Commit(ctx)
}

func (d *Auth) FindUserBySession(ctx context.Context, sessionID string) (*model.User, error) {
	user, err := model.Users.Query(
		model.SelectJoins.Users.InnerJoin.Sessions,
		model.SelectWhere.Sessions.ID.EQ(sessionID),
	).One(ctx, d.db)

	if err == nil {
		return user, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, app.ErrSessionNotFound
	}

	return user, oops.
		In("FindUserBySession").
		With("sessionID", sessionID).
		Errorf("cannot find user by session id: %w", err)
}
