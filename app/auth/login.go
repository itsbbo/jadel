package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type InsertSessionParam struct {
	UserID    ulid.ULID
	SessionID string
	Expires   time.Duration
	IPAddr    string
	UserAgent string
}

type FindByEmailPasswordQuery interface {
	FindByEmailPassword(ctx context.Context, email string, password string) (model.User, error)
}

type InsertSessionMutator interface {
	InsertSession(ctx context.Context, param InsertSessionParam) error
}

func (d *Deps) LoginPage(w http.ResponseWriter, r *http.Request) {
	d.server.RenderUI(w, r, "auth/login", app.NoUIProps)
}

func (d *Deps) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest

	if req, ok := d.server.Bind(w, r, loginSchema, &request); !ok {
		d.LoginPage(w, req)
		return
	}

	user, err := d.repo.FindByEmailPassword(r.Context(), request.Email, request.Password)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			d.LoginPage(w, d.server.AddValidationErrors(w, r, map[string]string{
				"email":    "Email or password is incorrect.",
				"password": "Email or password is incorrect.",
			}))
			return
		}

		slog.Error("Login.FindByEmailPassword failed", slog.Any("error", err))
		d.LoginPage(w, d.server.AddInternalErrorMsg(w, r))
		return
	}

	sessionID := rand.Text()
	err = d.repo.InsertSession(r.Context(), InsertSessionParam{
		UserID:    user.ID,
		SessionID: sessionID,
		Expires:   app.SessionTime,
		IPAddr:    r.RemoteAddr,
		UserAgent: r.UserAgent(),
	})
	if err != nil {
		slog.Error("Login.InsertSession failed", slog.Any("error", err))
		d.LoginPage(w, d.server.AddInternalErrorMsg(w, r))
		return
	}

	d.server.SetCookie(w, app.SessionKey, sessionID, app.SessionTime)
	d.server.RedirectTo(w, r, "/dashboard")
}
