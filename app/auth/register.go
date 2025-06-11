package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/itsbbo/jadel/app"
	"github.com/itsbbo/jadel/model"
	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	ConfirmationPassword string `json:"confirmationPassword"`
}

type NewUserWithSessionParam struct {
	Name      string
	Email     string
	Password  string
	IPAddr    string
	UserAgent string
}

type NewUserWithSessionMutator interface {
	NewUserWithSession(context.Context, NewUserWithSessionParam) (*model.User, string, error)
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

	_, session, err := d.repo.NewUserWithSession(c.Request().Context(), NewUserWithSessionParam{
		Name:      request.Name,
		Email:     request.Email,
		Password:  request.Password,
		IPAddr:    c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	})

	if err == nil {
		c.SetCookie(&http.Cookie{
			Name:    app.SessionKey,
			Value:   session,
			Expires: time.Now().Add(3 * time.Hour),
		})

		d.inertia.Redirect(c.Response(), c.Request(), "/dashboard")
		return nil
	}

	return err
}
