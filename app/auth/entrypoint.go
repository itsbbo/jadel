package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/itsbbo/jadel/app"
)

type Repository interface {
	NewUserWithSessionMutator
}

type Deps struct {
	server *app.Server
	repo   Repository
}

func New(server *app.Server, repo Repository) *Deps {
	return &Deps{
		server: server,
		repo:   repo,
	}
}

func (d *Deps) InitRoutes() {
	d.server.Route("/auth", func(r chi.Router) {
		r.Get("/register", d.RegisterPage)
		r.Post("/register", d.Register)
		r.Get("/login", d.LoginPage)
		r.Post("/login", d.Login)
	})
}
