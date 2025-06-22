package privatekey

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
)

type GetAllPrivateKeysQuery interface {
	GetAllPrivateKeys(ctx context.Context, userID ulid.ULID) ([]model.PrivateKey, error)
}

func (d *Deps) All(w http.ResponseWriter, r *http.Request) {
	user := app.CurrentUser(r)

	privateKeys, err := d.repo.GetAllPrivateKeys(r.Context(), user.ID)
	if err != nil {
		slog.Error("Cannot get all private keys", slog.Any("error", err))
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
		return
	}

	_ = json.NewEncoder(w).Encode(privateKeys)
}
