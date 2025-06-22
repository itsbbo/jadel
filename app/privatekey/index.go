package privatekey

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
)

type GetPrivateKeyIndexQuery interface {
	GetPrivateKeyIndex(ctx context.Context, param app.PaginationRequest) ([]model.PrivateKey, error)
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

	privateKeys, err := d.repo.GetPrivateKeyIndex(r.Context(), request)
	if err != nil {
		slog.Error("Cannot get private keys", slog.Any("error", err))
		d.server.Back(w, d.server.AddInternalErrorMsg(w, r))
		return
	}

	request.Limit = limit
	d.server.RenderUI(w, r, "privatekey/index", app.ToPaginationProps(request, privateKeys))
}
