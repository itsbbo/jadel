package settings

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
)

var (
	ErrEmailAlreadyTaken = errors.New("email already taken")
)

type ChangeProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DeleteAccountRequest struct {
	Password string `json:"password"`
}

type UpdateProfileParam struct {
	User  model.User
	Name  string
	Email string
}

type UpdateProfileMutator interface {
	UpdateProfile(ctx context.Context, param UpdateProfileParam) error
}

type DeleteAccountMutator interface {
	DeleteAccount(ctx context.Context, id ulid.ULID, password string) error
}

func (d *Deps) ProfilePage(w http.ResponseWriter, r *http.Request) {
	d.server.RenderUI(w, r, "settings/profile", app.NoUIProps)
}

func (d *Deps) ChangeProfile(w http.ResponseWriter, r *http.Request) {
	var request ChangeProfileRequest

	if req, ok := d.server.Bind(w, r, changeProfileSchema, &request); !ok {
		d.ProfilePage(w, req)
		return
	}

	user := app.CurrentUser(r)

	param := UpdateProfileParam{
		User:  user,
		Name:  request.Name,
		Email: request.Email,
	}

	err := d.repo.UpdateProfile(r.Context(), param)
	if err == nil {
		d.ProfilePage(w, r)
		return
	}

	if errors.Is(err, ErrEmailAlreadyTaken) {
		d.ProfilePage(w, d.server.AddValidationErrors(w, r, map[string]string{
			"email": "Email already taken",
		}))
		return
	}

	slog.Error("Error updating profile", slog.Any("error", err))
	d.ProfilePage(w, d.server.AddInternalErrorMsg(w, r))
}

func (d *Deps) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	var request DeleteAccountRequest

	if req, ok := d.server.Bind(w, r, destroyAccountSchema, &request); !ok {
		d.ProfilePage(w, req)
		return
	}

	user := app.CurrentUser(r)

	err := d.repo.DeleteAccount(r.Context(), user.ID, request.Password)
	if err == nil {
		d.server.SetCookie(w, app.SessionKey, "", -1)
		d.server.RedirectTo(w, r, "/auth/login")
	}

	if errors.Is(err, ErrWrongPassword) {
		d.ProfilePage(w, d.server.AddValidationErrors(w, r, map[string]string{
			"password": "Wrong password",
		}))
		return
	}

	slog.Error("Cannot update password", slog.Any("error", err))
	d.ProfilePage(w, d.server.AddInternalErrorMsg(w, r))
}
