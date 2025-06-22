package repo

import (
	"context"
	"slices"

	"github.com/guregu/null/v6"
	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/app/privatekey"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
	"github.com/samber/oops"
	"github.com/uptrace/bun"
)

type PrivateKey struct {
	db bun.IDB
}

func NewPrivateKey(db bun.IDB) *PrivateKey {
	return &PrivateKey{
		db: db,
	}
}

func (p *PrivateKey) GetPrivateKeyIndex(ctx context.Context, param app.PaginationRequest) ([]model.PrivateKey, error) {
	q := p.db.NewSelect().
		Column("id", "name", "description").
		Where("user_id = ?", param.UserID)

	switch {
	case !param.PrevID.IsZero():
		q = q.Where("id > ?", param.PrevID).Order("id ASC")
	case !param.NextID.IsZero():
		q = q.Where("id < ?", param.NextID).Order("id DESC")
	default:
		q = q.Order("id DESC")
	}

	var privateKeys []model.PrivateKey
	err := q.Model(&privateKeys).Limit(param.Limit).Scan(ctx, &privateKeys)
	if err != nil {
		return nil, oops.In("GetPrivateKeyIndex").With("param", param).Wrap(err)
	}

	if !param.PrevID.IsZero() {
		slices.Reverse(privateKeys)
	}

	return privateKeys, nil
}

func (p *PrivateKey) CreatePrivateKey(ctx context.Context, userID ulid.ULID, r privatekey.CreatePrivateKeyRequest) (model.PrivateKey, error) {
	privateKey := model.PrivateKey{
		ID:          ulid.Make(),
		UserID:      userID,
		Name:        r.Name,
		Description: null.StringFrom(r.Description),
		PublicKey:   r.PublicKey,
		PrivateKey:  r.PrivateKey,
	}

	_, err := p.db.NewInsert().
		Model(&privateKey).
		Exec(ctx)

	if err != nil {
		return model.PrivateKey{}, oops.In("CreatePrivateKey").With("param", r).Wrap(err)
	}

	return privateKey, nil
}
