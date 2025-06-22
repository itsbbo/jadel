package servers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
)

var (
	ErrUnknownPrivateKey = errors.New("unknown private key")
)

type CreateServerMutator interface {
	CreateServer(ctx context.Context, userID ulid.ULID, request CreateServerRequest) (model.Server, error)
}

type CreateServerRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	IP           string `json:"ip"`
	Port         int    `json:"port"`
	User         string `json:"user"`
	PrivateKeyID string `json:"privateKeyId"`
}

func (d *Deps) Create(w http.ResponseWriter, r *http.Request) {
	var request CreateServerRequest

	if req, ok := d.server.Bind(w, r, createServerSchema, &request); !ok {
		d.Index(w, req)
		return
	}

	user := app.CurrentUser(r)

	server, err := d.repo.CreateServer(r.Context(), user.ID, request)
	if err == nil {
		d.server.RedirectTo(w, r, fmt.Sprintf("/servers/%s", server.ID.String()))
		return
	}

	if errors.Is(err, ErrUnknownPrivateKey) {
		d.Index(w, d.server.AddValidationErrors(w, r, map[string]string{
			"private_key_id": "Unknown private key",
		}))
		return
	}

	slog.Error("Error creating server", slog.Any("error", err))
	d.Index(w, d.server.AddInternalErrorMsg(w, r))
}
