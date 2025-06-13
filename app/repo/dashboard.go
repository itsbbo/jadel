package repo

import (
	"context"

	"github.com/itsbbo/jadel/model"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/drivers/pgx"
)

type Dashboard struct {
	db pgx.Pool
}

func NewDashboard(db pgx.Pool) *Dashboard {
	return &Dashboard{db: db}
}

func (d *Dashboard) GetFiveLatestProjects(ctx context.Context) (model.ProjectSlice, error) {
	return model.Projects.Query(
		sm.Columns("id", "name"),
		sm.Limit(5),
		sm.OrderBy("id"),
	).All(ctx, d.db)
}

func (d *Dashboard) GetFiveLatestServers(ctx context.Context) (model.ServerSlice, error) {
	return model.Servers.Query(
		sm.Columns("id", "name"),
		sm.Limit(5),
		sm.OrderBy("id"),
	).All(ctx, d.db)
}
