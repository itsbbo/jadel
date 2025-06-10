package auth

import (
	"github.com/itsbbo/jadel/app"
	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	ConfirmationPassword string `json:"confirmationPassword"`
}

func (d *Deps) RegisterPage(c echo.Context) error {
	return d.inertia.Render(c.Response().Writer, c.Request(), "auth/register")
}

func (d *Deps) Register(c echo.Context) error {
	var request RegisterRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	if errs := registerSchema.Validate(&request); errs != nil {
		app.SetInertiaValidationErrorsZog(c, errs)
		return d.RegisterPage(c)
	}

	if request.Password != request.ConfirmationPassword {
		app.SetInertiaValidationErrorsMap(c, map[string]string{
			"passwordConfirmation": "Password input do not match the confirmation password.",
		})

		return d.RegisterPage(c)
	}

	return nil
}
