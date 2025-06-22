package repo

import (
	"context"
	"slices"

	"github.com/guregu/null/v6"
	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/app/servers"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
	"github.com/samber/oops"
	"github.com/uptrace/bun"
)

type Server struct {
	db bun.IDB
}

func NewServer(db bun.IDB) *Server {
	return &Server{db}
}

func (s *Server) GetServerIndex(ctx context.Context, param app.PaginationRequest) ([]model.Server, error) {
	q := s.db.NewSelect().
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

	var servers []model.Server
	err := q.Model(&servers).Limit(param.Limit).Scan(ctx, &servers)
	if err != nil {
		return nil, oops.In("GetServerIndex").With("param", param).Wrap(err)
	}

	if !param.PrevID.IsZero() {
		slices.Reverse(servers)
	}

	return servers, nil
}

func (s *Server) CreateServer(ctx context.Context, userID ulid.ULID, r servers.CreateServerRequest) (model.Server, error) {
	server := model.Server{
		ID:           ulid.Make(),
		UserID:       userID,
		Name:         r.Name,
		Description:  null.StringFrom(r.Description),
		IP:           r.IP,
		Port:         r.Port,
		PrivateKeyID: ulid.MustParse(r.PrivateKeyID),
	}

	_, err := s.db.NewInsert().
		Model(&server).
		Exec(ctx)

	if err == nil {
		return server, nil
	}

	if model.IsErrConstrainPrivateKeyInServer(err) {
		return model.Server{}, servers.ErrUnknownPrivateKey
	}

	return model.Server{}, oops.In("CreateServer").With("server", server).With("userID", userID.String()).Wrap(err)
}
