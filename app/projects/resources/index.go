package resources

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

type FindResourcesQuery interface {
	FindResources(ctx context.Context, userID, projectID, envID ulid.ULID) (model.Environment, error)
}

var (
	ErrEnvNotFound = errors.New("environment not found")
)

func (d *Deps) Index(w http.ResponseWriter, r *http.Request) {
	projectID, err := ulid.Parse(chi.URLParam(r, "project"))
	if err != nil {
		d.server.RenderNotFound(w, r)
		return
	}

	envID, err := ulid.Parse(chi.URLParam(r, "environment"))
	if err != nil {
		d.server.RenderNotFound(w, r)
		return
	}

	user := app.CurrentUser(r)

	env, err := d.repo.FindResources(r.Context(), user.ID, projectID, envID)
	if err == nil {
		d.server.RenderUI(w, r, "projects/resources/index", gonertia.Props{
			"env": env,
		})
		return
	}

	if errors.Is(err, ErrEnvNotFound) {
		d.server.RenderNotFound(w, r)
		return
	}

	slog.Error("cannot retrieved environments", slog.Any("error", err))
	d.server.Back(w, d.server.AddInternalErrorMsg(w, r))
}
