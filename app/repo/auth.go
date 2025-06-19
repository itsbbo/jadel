package repo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"time"

	"github.com/guregu/null/v6"
	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/app/auth"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
	"github.com/samber/oops"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	db bun.IDB
}

func NewAuth(db bun.IDB) *Auth {
	return &Auth{db: db}
}

func (d *Auth) NewUserWithSession(ctx context.Context, r auth.NewUserWithSessionParam) (model.User, string, error) {
	var (
		user      model.User
		sessionID = rand.Text()
	)

	err := d.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		user = model.User{
			ID:       ulid.Make(),
			Name:     r.Name,
			Email:    r.Email,
			Password: r.Password,
		}

		_, err := tx.NewInsert().Model(&user).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewInsert().Model(&model.Session{
			ID:        sessionID,
			UserID:    user.ID,
			IPAddress: null.StringFrom(r.IPAddr),
			UserAgent: null.StringFrom(r.UserAgent),
			ExpiredAt: time.Now().Add(app.SessionTime),
		}).Exec(ctx)

		return err
	})

	if err == nil {
		return user, sessionID, nil
	}

	if model.IsErrUniqueEmailUser(err) {
		return user, sessionID, auth.ErrDuplicateEmail
	}

	return user, "", oops.In("model.Users.Insert").Wrap(err)
}

func (d *Auth) FindByEmailPassword(ctx context.Context, email string, password string) (model.User, error) {
	var user model.User

	err := d.db.NewSelect().Model(&user).Where("email = ?", email).Limit(1).Scan(ctx)
	if err == nil {
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return user, auth.ErrUserNotFound
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return user, auth.ErrUserNotFound
	}

	return user, oops.
		In("FindByEmailPassword").
		With("email", email).
		Wrap(err)
}

func (d *Auth) InsertSession(ctx context.Context, param auth.InsertSessionParam) error {
	expiredAt := time.Now().Add(param.Expires)

	err := d.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(&model.Session{
			ID:        param.SessionID,
			UserID:    param.UserID,
			ExpiredAt: expiredAt,
			IPAddress: null.StringFrom(param.IPAddr),
			UserAgent: null.StringFrom(param.UserAgent),
		}).Exec(ctx)

		return err
	})

	if err == nil {
		return nil
	}

	return oops.
		In("InsertSession").
		With("userID", param.UserID.String()).
		With("sessionID", param.SessionID).
		Wrap(err)
}

func (d *Auth) FindUserBySession(ctx context.Context, sessionID string) (model.User, error) {
	var session model.Session

	err := d.db.NewSelect().Model(&session).
		Column("session.id").
		Where("session.id = ?", sessionID).
		Where("session.expired_at > ?", time.Now()).
		Limit(1).
		Relation("User").
		Scan(ctx)

	if err == nil {
		return session.User, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, app.ErrSessionNotFound
	}

	return model.User{}, oops.
		In("FindUserBySession").
		With("sessionID", sessionID).
		Wrap(err)
}

func (d *Auth) DeleteSession(ctx context.Context, sessionID string) error {
	err := d.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewDelete().
			Model(&model.Session{}).
			Where("sessions.id = ?", sessionID).
			Exec(ctx)

		return err
	})

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
