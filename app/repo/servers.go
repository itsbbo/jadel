package repo

import (
	"context"
	"slices"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
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
