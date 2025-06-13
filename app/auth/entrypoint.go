package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/itsbbo/jadel/app"
)

type Repository interface {
	NewUserWithSessionMutator
	FindByEmailPasswordQuery
	InsertSessionMutator
}

type Deps struct {
	server *app.Server
	middleware *app.Middleware
	repo   Repository
}

func New(server *app.Server, middleware *app.Middleware, repo Repository) *Deps {
	return &Deps{
		server: server,
		middleware: middleware,
		repo:   repo,
	}
}

func (d *Deps) InitRoutes() {
	d.server.Route("/auth", func(r chi.Router) {
		r.Use(d.middleware.RedirectIfAuthenticated)

		r.Get("/register", d.RegisterPage)
		r.Post("/register", d.Register)
		r.Get("/login", d.LoginPage)
		r.Post("/login", d.Login)
	})
}
