package dashboard

import (
	"context"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
	"github.com/romsar/gonertia/v2"
)

type Deps struct {
	server *app.Server
	repo   Repository
}

type Repository interface {
	GetFiveLatestProjects(ctx context.Context) (model.ProjectSlice, error)
	GetFiveLatestServers(ctx context.Context) (model.ServerSlice, error)
}

func New(server *app.Server, repo Repository) *Deps {
	return &Deps{server: server, repo: repo}
}

func (d *Deps) InitRoutes() {
	d.server.Get("/dashboard", d.DashboardUI)
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
