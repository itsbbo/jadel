package resources

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/itsbbo/jadel/app"
)

type Repository interface {
	FindResourcesQuery
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
	d.server.Route("/projects/{project}/environments/{environment}", func(r chi.Router) {
		r.Use(d.middleware.Auth)

		r.Get("/", d.Index)
		r.Get("/create", d.CreateUI)
	})
}

func (d *Deps) CreateUI(w http.ResponseWriter, r *http.Request) {
	d.server.RenderUI(w, r, "projects/resources/create", app.NoUIProps)
}
