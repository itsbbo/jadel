package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aarondl/opt/omit"
	"github.com/itsbbo/jadel/app/settings"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
	"github.com/samber/oops"
	"github.com/stephenafamo/bob/drivers/pgx"
	"golang.org/x/crypto/bcrypt"
)

type Settings struct {
	db pgx.Pool
}

func NewSettings(db pgx.Pool) *Settings {
	return &Settings{
		db: db,
	}
}

func (s *Settings) UpdatePassword(ctx context.Context, param settings.ChangePasswordParam) error {
	user, err := model.Users.Query(
		model.SelectWhere.Users.Email.EQ(param.Email),
	).One(ctx, s.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return settings.ErrUserNotFound
		}

		return oops.In("failed to find user").
			With("email", param.Email).
			Wrap(err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.CurrentPassword)); err != nil {
		return settings.ErrWrongPassword
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return oops.In("begin").With("email", param.Email).Wrap(err)
	}

	defer tx.Rollback(ctx)

	err = user.Update(ctx, tx, &model.UserSetter{
		Password: omit.From(param.NewPassword),
	})
	if err != nil {
		return oops.In("update password").With("email", param.Email).Wrap(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return oops.In("commit").With("email", param.Email).Wrap(err)
	}

	return nil
}

func (s *Settings) UpdateProfile(ctx context.Context, param settings.UpdateProfileParam) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return oops.In("begin").With("id", param.User.ID).Wrap(err)
	}

	defer tx.Rollback(ctx)

	err = param.User.Update(ctx, tx, &model.UserSetter{
		Name:  omit.From(param.Name),
		Email: omit.From(param.Email),
	})

	if err != nil {
		if model.UserErrors.ErrUniqueUsersEmailKey.Is(err) {
			return settings.ErrEmailAlreadyTaken
		}

		return oops.In("update profile").With("id", param.User.ID).Wrap(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return oops.In("commit").With("id", param.User.ID).Wrap(err)
	}

	return nil
}

func (s *Settings) DeleteAccount(ctx context.Context, id ulid.ULID, password string) error {
	user, err := model.FindUser(ctx, s.db, id, "password")

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return settings.ErrUserNotFound
		}

		return oops.In("DeleteAccount.FindUsersByID").
			With("id", id).
			Wrap(err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return settings.ErrWrongPassword
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return oops.In("begin").With("id", user.ID).Wrap(err)
	}

	defer tx.Rollback(ctx)

	err = user.Delete(ctx, tx)
	if err != nil {
		return oops.In("delete").With("id", user.ID).Wrap(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return oops.In("commit").With("id", user.ID).Wrap(err)
	}

	return nil
}
