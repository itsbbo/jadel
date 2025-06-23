package servers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/itsbbo/jadel/app"
)

type Repository interface {
	GetServerIndexQuery
	CreateServerMutator
}

type Deps struct {
	server     *app.Server
	repo       Repository
	middleware *app.Middleware
}

func New(server *app.Server, middleware *app.Middleware, repo Repository) *Deps {
	return &Deps{
		server:     server,
		middleware: middleware,
		repo:       repo,
	}
}

func (d *Deps) InitRoutes() {
	d.server.Route("/servers", func(r chi.Router) {
		r.Use(d.middleware.Auth)

		r.Get("/", d.Index)
		r.Post("/", d.Create)
		r.Get("/terminal", d.Experiment)
	})
}

func (d *Deps) Experiment(w http.ResponseWriter, r *http.Request) {
	d.server.RenderUI(w, r, "servers/terminal", app.NoUIProps)
}
