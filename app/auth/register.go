package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
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
	NewUserWithSession(context.Context, NewUserWithSessionParam) (model.User, string, error)
}

func (d *Deps) RegisterPage(w http.ResponseWriter, r *http.Request) {
	d.server.RenderUI(w, r, "auth/register", app.NoUIProps)
}

func (d *Deps) Register(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	if req, ok := d.server.Bind(w, r, registerSchema, &request); !ok {
		d.RegisterPage(w, req)
		return
	}

	if request.Password != request.PasswordConfirmation {
		d.RegisterPage(w, d.server.AddValidationErrors(w, r, map[string]string{
			"passwordConfirmation": "Password input do not match the confirmation password.",
		}))
		return
	}

	_, session, err := d.repo.NewUserWithSession(r.Context(), NewUserWithSessionParam{
		Name:      request.Name,
		Email:     request.Email,
		Password:  hashPassword(request.Password),
		IPAddr:    r.RemoteAddr,
		UserAgent: r.UserAgent(),
	})

	if err == nil {
		d.server.SetCookie(w, app.SessionKey, session, app.SessionTime)
		d.server.RedirectTo(w, r, "/dashboard")
		return
	}

	if errors.Is(err, ErrDuplicateEmail) {
		d.RegisterPage(w, d.server.AddValidationErrors(w, r, map[string]string{
			"email": "Email has already been taken.",
		}))
		return
	}

	slog.Error("register failed", slog.Any("error", err))
	d.RegisterPage(w, d.server.AddInternalErrorMsg(w, r))
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}
