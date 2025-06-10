package app

import (
	"net/http"

	"github.com/Oudwins/zog"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/romsar/gonertia/v2"
)

func NewServer() *echo.Echo {
	// static routes for javascript bundles
	e := echo.New()
	e.Static("/public/build/assets", "./public/build/assets")

	// non static routes
	g := e.Group("")
	g.Use(
		echomw.RemoveTrailingSlashWithConfig(echomw.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		echomw.Recover(),
		echomw.Secure(),
		echomw.RequestID(),
		echomw.CSRFWithConfig(echomw.CSRFConfig{
			TokenLookup:    "header:X-XSRF-TOKEN",
			CookieName:     "XSRF-TOKEN",
			CookiePath:     "/",
			CookieHTTPOnly: false,
			CookieSameSite: http.SameSiteStrictMode,
			ContextKey:     CSRFKey,
		}),
	)

	return e
}

func SetInertiaValidationErrorsZog(c echo.Context, errs zog.ZogIssueMap) {
	issues := make(gonertia.ValidationErrors, len(errs))

	for field, err := range errs {
		issues[field] = err[0].Message
	}

	ctx := gonertia.SetValidationErrors(c.Request().Context(), issues)
	c.SetRequest(c.Request().WithContext(ctx))
}

func SetInertiaValidationErrorsMap(c echo.Context, errs map[string]string) {
	issues := make(gonertia.ValidationErrors, len(errs))

	for field, err := range errs {
		issues[field] = err
	}

	ctx := gonertia.SetValidationErrors(c.Request().Context(), issues)
	c.SetRequest(c.Request().WithContext(ctx))
}