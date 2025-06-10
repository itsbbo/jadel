package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/romsar/gonertia/v2"
)

type Deps struct {
	inertia *gonertia.Inertia
	server  *echo.Echo
}

func New(inertia *gonertia.Inertia, server *echo.Echo) *Deps {
	return &Deps{
		inertia: inertia,
		server:  server,
	}
}

func (d *Deps) InitRoutes() {
	group := d.server.Group("/auth")

	group.GET("/register", d.RegisterPage)
	group.POST("/register", d.Register)
	group.GET("/login", d.LoginPage)
	group.POST("/login", d.Login)
}
