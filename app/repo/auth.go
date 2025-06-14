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
	"golang.org/x/crypto/bcrypt"
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
		if model.UserErrors.ErrUniqueUsersEmailKey.Is(err) {
			return nil, "", auth.ErrDuplicateEmail
		}

		return nil, "", oops.In("model.Users.Insert").Wrap(err)
	}

	session, err := model.Sessions.Insert(&model.SessionSetter{
		ID:        omit.From(rand.Text()),
		UserID:    omit.From(user.ID),
		IPAddress: omitnull.From(r.IPAddr),
		UserAgent: omitnull.From(r.UserAgent),
		ExpiredAt: omit.From(time.Now().Add(3 * time.Hour)),
	}).One(ctx, tx)

	if err != nil {
		return nil, "", oops.In("user.InsertSessions").Wrap(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, "", oops.In("tx.Commit").Wrap(err)
	}

	return user, session.ID, nil
}

func (d *Auth) FindByEmailPassword(ctx context.Context, email string, password string) (*model.User, error) {
	user, err := model.Users.Query(
		psql.WhereAnd(
			model.SelectWhere.Users.Email.EQ(email),
		),
	).One(ctx, d.db)

	if err == nil {
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return nil, auth.ErrUserNotFound
		}

		return user, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, auth.ErrUserNotFound
	}

	return user, oops.
		In("FindByEmailPassword").
		With("email", email).
		Wrap(err)
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
		Wrap(err)
}

func (d *Auth) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := model.Sessions.Delete(
		model.DeleteWhere.Sessions.ID.EQ(sessionID),
	).One(ctx, d.db)

	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return app.ErrSessionNotFound
	}

	return oops.
		In("DeleteSession").
		With("sessionID", sessionID).
		Wrap(err)
}
