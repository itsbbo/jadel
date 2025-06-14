package settings

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/itsbbo/jadel/app"
)

type Repository interface {
	UpdatePasswordMutator
	UpdateProfileMutator
	DeleteAccountMutator
}

type Deps struct {
	server     *app.Server
	middleware *app.Middleware
	repo       Repository
}

func New(server *app.Server, middleware *app.Middleware, repo Repository) *Deps {
	return &Deps{
		server:     server,
		middleware: middleware,
		repo:       repo,
	}
}

func (d *Deps) InitRoutes() {
	d.server.Route("/settings", func(r chi.Router) {
		r.Use(d.middleware.Auth)

		r.Get("/profile", d.ProfilePage)
		r.Patch("/profile", d.ChangeProfile)
		r.Delete("/profile", d.DeleteAccount)
		r.Get("/password", d.PasswordPage)
		r.Patch("/password", d.ChangePassword)
		r.Get("/appearance", d.Appearance)
	})
}

func (d *Deps) Appearance(w http.ResponseWriter, r *http.Request) {
	d.server.RenderUI(w, r, "settings/appearance", app.NoUIProps)
}
