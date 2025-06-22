package privatekey

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
)

type CreatePrivateKeyMutator interface {
	CreatePrivateKey(ctx context.Context, userID ulid.ULID, request CreatePrivateKeyRequest) (model.PrivateKey, error)
}

type CreatePrivateKeyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PublicKey   string `json:"publicKey"`
	PrivateKey  string `json:"privateKey"`
}

func (d *Deps) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePrivateKeyRequest

	if req, ok := d.server.Bind(w, r, createPrivateKeySchema, &req); !ok {
		d.Index(w, req)
		return
	}

	user := app.CurrentUser(r)

	privateKey, err := d.repo.CreatePrivateKey(r.Context(), user.ID, req)
	if err == nil {
		d.server.RedirectTo(w, r, fmt.Sprintf("/private-keys/%s", privateKey.ID.String()))
		return
	}

	slog.Error("Failed to create private key", slog.Any("error", err))
	d.Index(w, d.server.AddInternalErrorMsg(w, r))
}
