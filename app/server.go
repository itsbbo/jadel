package app

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Oudwins/zog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/itsbbo/jadel/gonertia"
)

var NoUIProps = gonertia.Props{}

type Server struct {
	*chi.Mux
	inertia        *gonertia.Inertia
	cookieSecure   bool
	cookieHTTPOnly bool
	cookieSameSite http.SameSite
	cookiePath     string
	cookieDomain   string
}

func (s *Server) AddStaticAssetsRoute() {
	publicAssets := http.StripPrefix(
		"/public/build/assets/",
		http.FileServer(http.Dir("./public/build/assets")),
	)

	s.Get("/public/build/assets/*", http.HandlerFunc(publicAssets.ServeHTTP))
}

func (s *Server) Bind(w http.ResponseWriter, r *http.Request, schema *zog.StructSchema, target any) (*http.Request, bool) {
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		slog.Error("cannot parse json", slog.Any("error", err))
		return nil, false
	}

	errMap := schema.Validate(target)
	if errMap == nil {
		return nil, true
	}

	issues := make(gonertia.ValidationErrors, len(errMap))
	for field, err := range errMap {
		issues[field] = err[0].Message
	}

	ctx := gonertia.SetValidationErrors(r.Context(), issues)
	return r.WithContext(ctx), false
}

func (s *Server) AddValidationErrors(_ http.ResponseWriter, r *http.Request, errMap map[string]string) *http.Request {
	issues := make(gonertia.ValidationErrors, len(errMap))
	for field, err := range errMap {
		issues[field] = err
	}

	ctx := gonertia.SetValidationErrors(r.Context(), issues)
	return r.WithContext(ctx)
}

func (s *Server) AddInternalErrorMsg(w http.ResponseWriter, r *http.Request) *http.Request {
	return s.AddValidationErrors(w, r, map[string]string{
		"_internal": "Internal server error. Try again later.",
	})
}

func (s *Server) SetCookie(w http.ResponseWriter, key, val string, expiry time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    val,
		Secure:   s.cookieSecure,
		SameSite: s.cookieSameSite,
		Path:     s.cookiePath,
		Domain:   s.cookieDomain,
		HttpOnly: s.cookieHTTPOnly,
		MaxAge:   int(expiry),
	})
}

func (s *Server) RenderUI(w http.ResponseWriter, r *http.Request, components string, props gonertia.Props) {
	_ = s.inertia.Render(w, r, components, props)
}

func (s *Server) RedirectTo(w http.ResponseWriter, r *http.Request, path string) {
	s.inertia.Redirect(w, r, path)
}

func (s *Server) Back(w http.ResponseWriter, r *http.Request) {
	s.inertia.Back(w, r)
}

func (s *Server) RenderNotFound(w http.ResponseWriter, r *http.Request) {
	s.RenderUI(w, r, "error/not-found", NoUIProps)
}

func (s *Server) PrintRoutes() {
	chi.Walk(s.Mux, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		slog.Info("Route", slog.String("method", method), slog.String("route", route))
		return nil
	})
}

func NewServer(c Config, inertia *gonertia.Inertia) *Server {
	r := chi.NewRouter()

	var sameSite http.SameSite
	switch c.Cookie.SameSite {
	case "Lax", "lax":
		sameSite = http.SameSiteLaxMode
	case "Strict", "strict":
		sameSite = http.SameSiteStrictMode
	case "None", "none":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteDefaultMode
	}

	r.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
	)

	server := Server{
		Mux:            r,
		inertia:        inertia,
		cookieSecure:   c.Cookie.Secure,
		cookieHTTPOnly: c.Cookie.HTTPOnly,
		cookiePath:     c.Cookie.Path,
		cookieDomain:   c.Cookie.Domain,
		cookieSameSite: sameSite,
	}

	server.NotFound(func(w http.ResponseWriter, r *http.Request) {
		server.RenderUI(w, r, "error/not-found", NoUIProps)
	})

	server.AddStaticAssetsRoute()

	return &server
}
