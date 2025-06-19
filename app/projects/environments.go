package projects

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/gonertia"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
)

var (
	ErrProjectNotFound = errors.New("project not found")
)

type AllEnvironmentsQuery interface {
	AllEnvironments(ctx context.Context, userID, projectID ulid.ULID) (model.Project, error)
}

func (d *Deps) Environments(w http.ResponseWriter, r *http.Request) {
	userID := app.CurrentUser(r).ID

	projectID, err := ulid.Parse(chi.URLParam(r, "project"))
	if err != nil {
		d.server.RenderNotFound(w, r)
		return
	}

	project, err := d.repo.AllEnvironments(r.Context(), userID, projectID)
	if err == nil {
		d.server.RenderUI(w, r, "projects/environments", gonertia.Props{
			"project": project,
			"envs":    project.Environments,
		})
		return
	}

	if errors.Is(err, ErrProjectNotFound) {
		d.server.RenderNotFound(w, r)
		return
	}

	slog.Error("cannot get environments", slog.Any("error", err))
	d.server.Back(w, d.server.AddInternalErrorMsg(w, r))
}
