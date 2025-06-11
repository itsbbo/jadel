package auth

import (
	"net/http"

	"github.com/itsbbo/jadel/app"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d *Deps) LoginPage(w http.ResponseWriter, r *http.Request) {
	d.server.RenderUI(w, r, "auth/login", app.NoUIProps)
}

func (d *Deps) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest

	ok := d.server.BindJSON(w, r, loginSchema, &request)
	if !ok {
		return
	}
}
