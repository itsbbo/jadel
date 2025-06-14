package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/app"
)

type DeleteSessionMutator interface {
	DeleteSession(ctx context.Context, sessionID string) error
}

func (d *Deps) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(app.SessionKey)
	if err != nil {
		d.server.Back(w, r)
		return
	}

	err = d.repo.DeleteSession(r.Context(), cookie.Value)
	if err == nil || errors.Is(err, app.ErrSessionNotFound) {
		d.server.SetCookie(w, app.SessionKey, "", -1)
		d.server.RedirectTo(w, r, "/auth/login")
		return
	}

	slog.Error("logout failed", slog.Any("error", err))
	d.server.Back(w, d.server.AddInternalErrorMsg(w, r))
}
