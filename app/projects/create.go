package projects

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
)

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateProjectParam struct {
	User        *model.User
	Name        string
	Description string
}

type CreateProjectMutator interface {
	CreateProject(ctx context.Context, param CreateProjectParam) (*model.Project, error)
}

func (d *Deps) CreateProject(w http.ResponseWriter, r *http.Request) {
	var request CreateProjectRequest

	if req, ok := d.server.Bind(w, r, createProjectSchema, &request); !ok {
		d.Index(w, req)
		return
	}

	user := app.CurrentUser(r)

	param := CreateProjectParam{
		User:        user,
		Name:        request.Name,
		Description: request.Description,
	}

	project, err := d.repo.CreateProject(r.Context(), param)
	if err != nil {
		slog.Error("Error creating project", slog.Any("error", err))
		d.Index(w, d.server.AddInternalErrorMsg(w, r))
		return
	}

	d.server.RedirectTo(w, r, "/projects/"+project.ID.String())
}
