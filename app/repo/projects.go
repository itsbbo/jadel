package repo

import (
	"context"

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

	if !param.NextID.IsZero() {
		q.Apply(
			model.SelectWhere.Projects.ID.LT(param.NextID),
			sm.OrderBy(model.ColumnNames.Projects.ID).Desc(),
		)
		return q.All(ctx, p.db)
	}

	if !param.PrevID.IsZero() {
		q.Apply(
			model.SelectWhere.Projects.ID.GT(param.PrevID),
			sm.OrderBy(model.ColumnNames.Projects.ID).Asc(),
		)
		return q.All(ctx, p.db)
	}

	q.Apply(sm.OrderBy(model.ColumnNames.Projects.ID).Desc())
	return q.All(ctx, p.db)
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
