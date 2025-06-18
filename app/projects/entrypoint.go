package projects

import (
	"github.com/go-chi/chi/v5"
	"github.com/itsbbo/jadel/app"
)

type Repository interface {
	GetProjectIndexQuery
	CreateProjectMutator
	AllEnvironmentsQuery
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
	d.server.Route("/projects", func(r chi.Router) {
		r.Use(d.middleware.Auth)

		r.Get("/", d.Index)
		r.Post("/", d.CreateProject)
		r.Get("/{project}/environments", d.Environments)
	})
}
