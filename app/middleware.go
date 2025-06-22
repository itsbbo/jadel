package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/itsbbo/jadel/gonertia"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrEnvNotFound     = errors.New("env not found")
)

type FindSessionQuery interface {
	FindUserBySession(ctx context.Context, sessionID string) (model.User, error)
}

type FindResourcesGeneralQuery interface {
	FindResourcesGeneral(ctx context.Context, userID, projectID, envID ulid.ULID) (model.Environment, error)
}

type Middleware struct {
	server *Server
	auth   FindSessionQuery
	env    FindResourcesGeneralQuery
}

func NewMiddleware(server *Server, auth FindSessionQuery, env FindResourcesGeneralQuery) *Middleware {
	return &Middleware{
		server: server,
		auth:   auth,
		env:    env,
	}
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie(SessionKey)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				m.server.RedirectTo(w, r, "/auth/login")
				return
			}

			slog.Error("Cannot get session cookie", slog.Any("error", err))
			m.server.RedirectTo(w, r, "/auth/login")
			return
		}

		user, err := m.auth.FindUserBySession(r.Context(), sessionID.Value)
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
	})
}

func (m *Middleware) LoadResources(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectID, err := ulid.Parse(chi.URLParam(r, "project"))
		if err != nil {
			m.server.RenderNotFound(w, r)
			return
		}

		envID, err := ulid.Parse(chi.URLParam(r, "environment"))
		if err != nil {
			m.server.RenderNotFound(w, r)
			return
		}

		user := CurrentUser(r)

		env, err := m.env.FindResourcesGeneral(r.Context(), user.ID, projectID, envID)
		if err == nil {
			ctx := gonertia.SetProp(r.Context(), EnvKey, env)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if errors.Is(err, ErrEnvNotFound) {
			m.server.RenderNotFound(w, r)
			return
		}

		slog.Error("cannot retrieved environments", slog.Any("error", err))
		m.server.Back(w, m.server.AddInternalErrorMsg(w, r))
	})
}
