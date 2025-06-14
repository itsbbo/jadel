package settings

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
)

type ChangePasswordRequest struct {
	CurrentPassword      string `json:"currentPassword"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

type ChangePasswordParam struct {
	Email           string
	CurrentPassword string
	NewPassword     string
}

type UpdatePasswordMutator interface {
	UpdatePassword(ctx context.Context, param ChangePasswordParam) error
}

func (d *Deps) PasswordPage(w http.ResponseWriter, r *http.Request) {
	d.server.RenderUI(w, r, "settings/password", app.NoUIProps)
}

func (d *Deps) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var request ChangePasswordRequest

	if req, ok := d.server.Bind(w, r, changePasswordSchema, &request); !ok {
		d.PasswordPage(w, req)
		return
	}

	if request.Password != request.PasswordConfirmation {
		d.PasswordPage(w, d.server.AddValidationErrors(w, r, map[string]string{
			"password":             "Passwords do not match",
			"passwordConfirmation": "Passwords do not match",
		}))

		return
	}

	user := app.CurrentUser(r)

	param := ChangePasswordParam{
		Email:           user.Email,
		CurrentPassword: request.CurrentPassword,
		NewPassword:     hashPassword(request.Password),
	}

	err := d.repo.UpdatePassword(r.Context(), param)
	if err == nil {
		d.PasswordPage(w, r)
		return
	}

	if errors.Is(err, ErrUserNotFound) {
		d.PasswordPage(w, d.server.AddValidationErrors(w, r, map[string]string{
			"currentPassword": "User not found",
		}))
		return
	}

	if errors.Is(err, ErrWrongPassword) {
		d.PasswordPage(w, d.server.AddValidationErrors(w, r, map[string]string{
			"currentPassword": "Incorrect password",
		}))
		return
	}

	slog.Error("Error updating password", slog.Any("error", err))
	d.PasswordPage(w, d.server.AddInternalErrorMsg(w, r))
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}
