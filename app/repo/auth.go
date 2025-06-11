package repo

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/itsbbo/jadel/app/auth"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
	"github.com/stephenafamo/bob"
)

type Auth struct {
	db bob.DB
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

	sessionsId := rand.Text()
	sessionsExpired := time.Now().Add(3 * time.Hour)

	user.InsertSessions(ctx, d.db, &model.SessionSetter{
		ID:        &sessionsId,
		UserID:    &userId,
		ExpiredAt: &sessionsExpired,
	})

	err = tx.Commit(ctx)
	if err != nil {
		return nil, "", err
	}

	return user, sessionsId, nil
}
