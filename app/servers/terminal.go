package servers

import (
	"net/http"

	"github.com/itsbbo/jadel/app"
)

func (d *Deps) Terminal(w http.ResponseWriter, r *http.Request) {
	conn, err := app.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		d.server.RenderUI(w, d.server.AddInternalErrorMsg(w, r), "servers/terminal", app.NoUIProps)
		return
	}

	defer conn.Close()

	d.server.RenderUI(w, r, "servers/terminal", app.NoUIProps)
}
