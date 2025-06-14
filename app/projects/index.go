package projects

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
)

type GetProjectIndexQuery interface {
	GetProjectIndex(ctx context.Context, param IndexRequest) (model.ProjectSlice, error)
}

type IndexRequest struct {
	PaginationMode app.PaginationMode
	Limit          int
	TrackID        ulid.ULID
	UserID         ulid.ULID
}

func (d *Deps) Index(w http.ResponseWriter, r *http.Request) {
	trackID, err := ulid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		trackID = ulid.Zero
	}

	user := app.CurrentUser(r)

	projects, err := d.repo.GetProjectIndex(r.Context(), IndexRequest{
		PaginationMode: app.ToPaginationMode(r.URL.Query().Get("mode")),
		Limit:          app.PaginationDefaultLimit,
		TrackID:        trackID,
		UserID:         user.ID,
	})

	if err != nil {
		slog.Error("Cannot get project", slog.Any("error", err))
		d.server.Back(w, d.server.AddInternalErrorMsg(w, r))
	}

	d.server.RenderUI(
		w, r, "projects/index",
		app.ToPaginationProps("/projects", app.PaginationDefaultLimit, projects),
	)
}
