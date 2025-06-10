package auth

import "github.com/labstack/echo/v4"

func (d *Deps) LoginUI(c echo.Context) error {
	return d.inertia.Render(c.Response().Writer, c.Request(), "auth/login")
}
