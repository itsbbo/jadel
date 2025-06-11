package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
)

type RegisterRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

type NewUserWithSessionParam struct {
	Name      string
	Email     string
	Password  string
	IPAddr    string
	UserAgent string
}

type NewUserWithSessionMutator interface {
	NewUserWithSession(context.Context, NewUserWithSessionParam) (*model.User, string, error)
}

func (d *Deps) RegisterPage(w http.ResponseWriter, r *http.Request) {
	d.server.RenderUI(w, r, "auth/register", app.NoUIProps)
}

func (d *Deps) Register(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	ok := d.server.BindJSON(w, r, registerSchema, &request)
	if !ok {
		return
	}

	if request.Password != request.PasswordConfirmation {
		d.server.AddValidationErrors(w, r, map[string]string{
			"passwordConfirmation": "Password input do not match the confirmation password.",
		})
		d.RegisterPage(w, r)
		return
	}

	_, session, err := d.repo.NewUserWithSession(r.Context(), NewUserWithSessionParam{
		Name:      request.Name,
		Email:     request.Email,
		Password:  request.Password,
		IPAddr:    d.server.RealIP(r),
		UserAgent: r.UserAgent(),
	})

	if err == nil {
		d.server.SetCookie(w, app.SessionKey, session, 3*time.Hour)
		d.server.RedirectTo(w, r, "/dashboard")
		return
	}

	if errors.Is(err, model.ErrUniqueConstraint) {
		d.server.AddValidationErrors(w, r, map[string]string{
			"email": "Email has already been taken.",
		})
		d.RegisterPage(w, r)
		return
	}

	slog.Error("register failed", slog.Any("error", err))
	d.server.AddInternalErrorMsg(w, r)
	d.RegisterPage(w, r)
}
