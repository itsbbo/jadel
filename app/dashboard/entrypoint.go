package dashboard

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/gonertia"
	"github.com/itsbbo/jadel/model"
)

type Deps struct {
	server     *app.Server
	middleware *app.Middleware
	repo       Repository
}

type Repository interface {
	GetFiveLatestProjects(ctx context.Context) (model.ProjectSlice, error)
	GetFiveLatestServers(ctx context.Context) (model.ServerSlice, error)
}

func New(server *app.Server, middleware *app.Middleware, repo Repository) *Deps {
	return &Deps{server: server, middleware: middleware, repo: repo}
}

func (d *Deps) InitRoutes() {
	d.server.Group(func(r chi.Router) {
		r.Use(d.middleware.Auth)
		r.Get("/dashboard", d.DashboardUI)
	})
}

func (d *Deps) DashboardUI(w http.ResponseWriter, r *http.Request) {
	projects, err := d.repo.GetFiveLatestProjects(r.Context())
	if err != nil {
		d.server.AddInternalErrorMsg(w, r)
		d.server.RenderUI(w, r, "dashboard", app.NoUIProps)
		return
	}

	servers, err := d.repo.GetFiveLatestServers(r.Context())
	if err != nil {
		d.server.AddInternalErrorMsg(w, r)
		d.server.RenderUI(w, r, "dashboard", app.NoUIProps)
		return
	}

	d.server.RenderUI(w, r, "dashboard", gonertia.Props{
		"projects": projects,
		"servers":  servers,
	})
}
