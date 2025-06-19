package repo

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"github.com/guregu/null/v6"
	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/app/projects"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
	"github.com/samber/oops"
	"github.com/uptrace/bun"
)

type Project struct {
	db bun.IDB
}

func NewProject(db bun.IDB) *Project {
	return &Project{
		db: db,
	}
}

func (p *Project) GetProjectIndex(ctx context.Context, param app.PaginationRequest) ([]model.Project, error) {
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

	var projects []model.Project
	err := q.Model(&projects).Limit(param.Limit).Scan(ctx, &projects)
	if err != nil {
		return nil, oops.In("GetProjectIndex").With("param", param).Wrap(err)
	}

	if !param.PrevID.IsZero() {
		slices.Reverse(projects)
	}

	return projects, nil
}

func (p *Project) CreateProject(ctx context.Context, param projects.CreateProjectParam) (model.Project, error) {
	project := model.Project{
		ID:          ulid.Make(),
		Name:        param.Name,
		UserID:      param.User.ID,
		Description: null.StringFrom(param.Description),
	}

	err := p.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(&project).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewInsert().Model(&model.Environment{
			ID:        ulid.Make(),
			ProjectID: project.ID,
			Name:      "production",
		}).Exec(ctx)

		return err
	})

	if err != nil {
		return model.Project{}, oops.
			In("CreateProject").
			With("name", param.Name).
			With("description", param.Description).
			With("userID", param.User.ID.String()).
			Wrap(err)
	}

	return project, nil
}

func (p *Project) AllEnvironments(ctx context.Context, userID, projectID ulid.ULID) (model.Project, error) {
	var project model.Project

	err := p.db.NewSelect().
		Model(&project).
		Column("id", "name").
		Relation("Environments").
		Where("id = ?", projectID).
		Where("user_id = ?", userID).
		Order("id DESC").
		Limit(10).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Project{}, projects.ErrProjectNotFound
		}

		return model.Project{}, oops.In("ProjectQuery").With("projectID", projectID.String()).Wrap(err)
	}

	return project, nil
}

func (p *Project) FindSpesificEnvironments(ctx context.Context, userID, projectID, envID ulid.ULID) (model.Environment, error) {
	var env model.Environment

	err := p.db.NewSelect().
		Model(&env).
		Column("id", "name").
		Relation("Project", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "name").
				Where("project.id = ?", projectID).
				Where("project.user_id = ?", userID)
		}).
		Where("environment.id = ?", envID).
		Scan(ctx)

	if err == nil {
		return env, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return env, app.ErrEnvNotFound
	}

	return env, oops.
		In("project.LoadEnvironments").
		With("userID", userID.String()).
		With("envID", envID.String()).
		Wrap(err)
}
