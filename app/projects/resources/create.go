package resources

import (
	"net/http"

	"github.com/itsbbo/jadel/app"
)

func (d *Deps) CreateUI(w http.ResponseWriter, r *http.Request) {
	d.server.RenderUI(w, r, "projects/resources/create", app.NoUIProps)
}
