package auth

import (
	"github.com/itsbbo/jadel/app"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d *Deps) LoginPage(c echo.Context) error {
	return d.inertia.Render(c.Response().Writer, c.Request(), "auth/login")
}

func (d *Deps) Login(c echo.Context) error {
	var request LoginRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	if errs := loginSchema.Validate(&request); errs != nil {
		app.SetInertiaValidationErrorsZog(c, errs)
		return d.LoginPage(c)
	}

	return nil
}
