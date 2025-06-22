package resources

import (
	"context"
	"errors"
	"net/http"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/gonertia"
	"github.com/itsbbo/jadel/model"
	"github.com/oklog/ulid/v2"
)

type FindResourcesGeneralQuery interface {
	FindResourcesGeneral(ctx context.Context, userID, projectID, envID ulid.ULID) (model.Environment, error)
}

var (
	ErrEnvNotFound = errors.New("environment not found")
)

func (d *Deps) Index(w http.ResponseWriter, r *http.Request) {
	env := app.CurrentEnv(r)
	d.server.RenderUI(w, r, "projects/resources/index", gonertia.Props{
		"env": env,
	})
}
