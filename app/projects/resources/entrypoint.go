package resources

import (
	"github.com/go-chi/chi/v5"
	"github.com/itsbbo/jadel/app"
)

type Deps struct {
	server     *app.Server
	middleware *app.Middleware
}

func New(server *app.Server, middleware *app.Middleware) *Deps {
	return &Deps{
		server:     server,
		middleware: middleware,
	}
}

func (d *Deps) InitRoutes() {
	d.server.Route("/projects/{project}/environments/{environment}", func(r chi.Router) {
		r.Use(d.middleware.Auth)
		r.Get("/create", d.CreateUI)
	})
}
