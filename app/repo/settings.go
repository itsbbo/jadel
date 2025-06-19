package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/itsbbo/jadel/app/settings"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
	"github.com/samber/oops"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type Settings struct {
	db bun.IDB
}

func NewSettings(db bun.IDB) *Settings {
	return &Settings{
		db: db,
	}
}

func (s *Settings) UpdatePassword(ctx context.Context, param settings.ChangePasswordParam) error {
	var user model.User

	err := s.db.NewSelect().
		Model(&user).
		Where("email = ?", param.Email).
		Scan(ctx)

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

	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err = tx.NewUpdate().
			Model(&user).
			Set("password = ?", param.NewPassword).
			WherePK().
			Exec(ctx)

		return err
	})

	if err != nil {
		return oops.In("update password").With("email", param.Email).Wrap(err)
	}

	return nil
}

func (s *Settings) UpdateProfile(ctx context.Context, param settings.UpdateProfileParam) error {
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewUpdate().
			Model(&param.User).
			Set("name = ?", param.Name).
			Set("email = ?", param.Email).
			WherePK().
			Exec(ctx)

		return err
	})

	if err == nil {
		return nil
	}

	if model.IsErrUniqueEmailUser(err) {
		return settings.ErrEmailAlreadyTaken
	}

	return oops.In("update profile").With("id", param.User.ID).Wrap(err)
}

func (s *Settings) DeleteAccount(ctx context.Context, id ulid.ULID, password string) error {
	var user model.User

	err := s.db.NewSelect().
		Model(&user).
		Where("id = ?", id).
		Scan(ctx)

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

	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, errTx := tx.NewDelete().
			Model(&user).
			WherePK().
			Exec(ctx)

		return errTx
	})

	if err != nil {
		return oops.In("delete").With("id", user.ID).Wrap(err)
	}

	return nil
}
