package servers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
)

type GetServerIndexQuery interface {
	GetServerIndex(ctx context.Context, param app.PaginationRequest) ([]model.Server, error)
}

func (d *Deps) Index(w http.ResponseWriter, r *http.Request) {
	prevId, err := ulid.Parse(r.URL.Query().Get("prevId"))
	if err != nil {
		prevId = ulid.Zero
	}

	nextId, err := ulid.Parse(r.URL.Query().Get("nextId"))
	if err != nil {
		nextId = ulid.Zero
	}

	if !prevId.IsZero() && !nextId.IsZero() {
		nextId, prevId = ulid.Zero, ulid.Zero
	}

	limit := app.PaginationDefaultLimit

	user := app.CurrentUser(r)
	request := app.PaginationRequest{
		PrevID: prevId,
		NextID: nextId,
		Limit:  limit + 1,
		UserID: user.ID,
	}

	servers, err := d.repo.GetServerIndex(r.Context(), request)
	if err != nil {
		slog.Error("Cannot get server", slog.Any("error", err))
		d.server.Back(w, d.server.AddInternalErrorMsg(w, r))
	}

	request.Limit = limit
	d.server.RenderUI(w, r, "servers/index", app.ToPaginationProps(request, servers))
}
