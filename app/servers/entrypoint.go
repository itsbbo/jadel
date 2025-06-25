package servers

import (
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

		d.server.Route("/{server}", func(rr chi.Router) {
			rr.Get("/configuration", nil)
			rr.Get("/proxy", nil)
			rr.Get("/resources", nil)
			rr.Get("/terminal", d.Terminal)
			rr.Get("/security", nil)
		})

	})
}
