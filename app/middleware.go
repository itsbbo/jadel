package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/gonertia"
	"github.com/itsbbo/jadel/model"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type Repository interface {
	FindUserBySession(ctx context.Context, sessionID string) (*model.User, error)
}

type Middleware struct {
	server *Server
	repo   Repository
}

func NewMiddleware(server *Server, repo Repository) *Middleware {
	return &Middleware{
		server: server,
		repo:   repo,
	}
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie(SessionKey)
		if err != nil {
			m.server.RedirectTo(w, r, "/auth/login")
			return
		}

		user, err := m.repo.FindUserBySession(r.Context(), sessionID.Value)
		if err == nil {
			ctx := gonertia.SetProp(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if errors.Is(err, ErrSessionNotFound) {
			m.server.SetCookie(w, SessionKey, "", -1)
			m.server.RedirectTo(w, r, "/auth/login")
			return
		}

		slog.Error("Cannot check auth user", slog.Any("error", err))
		m.server.RenderUI(w, r, "error/internal-server-error", NoUIProps)
		return
	})
}

func (m *Middleware) RedirectIfAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SessionKey)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if cookie.Value == "" {
			m.server.SetCookie(w, SessionKey, "", -1)
			m.server.RedirectTo(w, r, "/auth/login")
			return
		}

		m.server.RedirectTo(w, r, "/dashboard")
		return
	})
}
