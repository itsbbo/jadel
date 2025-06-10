package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
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