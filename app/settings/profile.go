package settings

import (
	"net/http"

	"github.com/itsbbo/jadel/app"
)

func (d *Deps) ProfilePage(w http.ResponseWriter, r *http.Request) {
	d.server.RenderUI(w, r, "settings/profile", app.NoUIProps)
}