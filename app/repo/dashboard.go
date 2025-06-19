package repo

import (
	"context"

	"github.com/itsbbo/jadel/model"
	"github.com/samber/oops"
	"github.com/uptrace/bun"
)

type Dashboard struct {
	db bun.IDB
}

func NewDashboard(db bun.IDB) *Dashboard {
	return &Dashboard{db: db}
}

func (d *Dashboard) GetFiveLatestProjects(ctx context.Context) ([]model.Project, error) {
	var projects []model.Project

	err := d.db.NewSelect().
		Model(&projects).
		Column("id", "name").
		Limit(5).
		Order("id DESC").
		Scan(ctx)

	if err != nil {
		return nil, oops.In("GetFiveLatestProjects").Wrap(err)
	}

	return projects, nil
}

func (d *Dashboard) GetFiveLatestServers(ctx context.Context) ([]model.Server, error) {
	var servers []model.Server

	err := d.db.NewSelect().
		Model(&servers).
		Column("id", "name").
		Limit(5).
		Order("id DESC").
		Scan(ctx)

	if err != nil {
		return nil, oops.In("GetFiveLatestServers").Wrap(err)
	}

	return servers, nil
}
