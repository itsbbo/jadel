package repo

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/app/projects"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
	"github.com/samber/oops"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/drivers/pgx"
)

type Project struct {
	db pgx.Pool
}

func NewProject(db pgx.Pool) *Project {
	return &Project{
		db: db,
	}
}

func (p *Project) GetProjectIndex(ctx context.Context, param app.PaginationRequest) (model.ProjectSlice, error) {
	q := model.Projects.Query(
		sm.Columns(
			model.ColumnNames.Projects.ID,
			model.ColumnNames.Projects.Name,
			model.ColumnNames.Projects.Description,
		),
		model.SelectWhere.Projects.UserID.EQ(param.UserID),
		sm.Limit(param.Limit),
	)

	switch {
	case !param.PrevID.IsZero():
		q.Apply(
			model.SelectWhere.Projects.ID.GT(param.PrevID),
			sm.OrderBy(model.ColumnNames.Projects.ID).Asc(),
		)

	case !param.NextID.IsZero():
		q.Apply(
			model.SelectWhere.Projects.ID.LT(param.NextID),
			sm.OrderBy(model.ColumnNames.Projects.ID).Desc(),
		)

	default:
		q.Apply(sm.OrderBy(model.ColumnNames.Projects.ID).Desc())
	}

	results, err := q.All(ctx, p.db)
	if err != nil {
		return nil, oops.In("GetProjectIndex").With("param", param).Wrap(err)
	}

	if !param.PrevID.IsZero() {
		slices.Reverse(results)
	}

	return results, nil
}

func (p *Project) CreateProject(ctx context.Context, param projects.CreateProjectParam) (*model.Project, error) {
	projectSetter := &model.ProjectSetter{
		ID:          omit.From(ulid.Make()),
		Name:        omit.From(param.Name),
		UserID:      omit.From(param.User.ID),
		Description: omitnull.From(param.Description),
	}

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, oops.In("Begin").Wrap(err)
	}

	defer tx.Rollback(ctx)

	errWrap := oops.
		With("name", param.Name).
		With("description", param.Description).
		With("userID", param.User.ID.String())

	project, err := model.Projects.Insert(projectSetter).One(ctx, tx)
	if err != nil {
		return nil, errWrap.In("Projects.Insert").Wrap(err)
	}

	err = project.InsertEnvironments(ctx, tx, &model.EnvironmentSetter{
		ID:   omit.From(ulid.Make()),
		Name: omit.From("production"),
	})

	if err != nil {
		return nil, errWrap.In("InsertEnvironments").Wrap(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errWrap.In("tx.Commit").Wrap(err)
	}

	return project, nil
}

func (p *Project) AllEnvironments(ctx context.Context, userID, projectID ulid.ULID) (*model.Project, model.EnvironmentSlice, error) {
	project, err := model.Projects.Query(
		sm.Columns(model.ColumnNames.Projects.ID, model.ColumnNames.Projects.Name),
		model.SelectWhere.Projects.UserID.EQ(userID),
		model.SelectWhere.Projects.ID.EQ(projectID),
	).One(ctx, p.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}

		return nil, nil, oops.In("ProjectQuery").With("projectID", projectID.String()).Wrap(err)
	}

	environments, err := project.Environments(
		sm.Columns(model.ColumnNames.Environments.ID, model.ColumnNames.Environments.Name),
	).All(ctx, p.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return project, nil, nil
		}

		return nil, nil, oops.In("ProjectEnvironments").With("projectID", projectID.String()).Wrap(err)
	}

	return project, environments, nil
}

func (p *Project) FindSpesificEnvironments(ctx context.Context, userID, projectID, envID ulid.ULID) (*model.Environment, error) {
	project, err := model.Projects.Query().One(ctx, p.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.
				In("FindSpesificEnvironments - model.Projects.Query()").
				With("userID", userID.String()).
				With("envID", envID.String()).
				Wrap(err)
		}
	}

	err = project.LoadEnvironments(ctx, p.db, model.SelectWhere.Environments.ID.EQ(envID))
	if err == nil {
		env := *project.R.Environments[0]
		copyProject := *project
		copyProject.R.Environments = nil
		env.R.Project = &copyProject
		return &env, nil
	}

	return nil, oops.
		In("project.LoadEnvironments").
		With("userID", userID.String()).
		With("envID", envID.String()).
		Wrap(err)
}
