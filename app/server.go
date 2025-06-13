package app

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Oudwins/zog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/romsar/gonertia/v2"
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

func (s *Server) BindJSON(_ http.ResponseWriter, r *http.Request, schema *zog.StructSchema, target any) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return false
	}

	errMap := schema.Validate(target)
	if errMap == nil {
		return true
	}

	issues := make(gonertia.ValidationErrors, len(errMap))
	for field, err := range errMap {
		issues[field] = err[0].Message
	}

	ctx := gonertia.SetValidationErrors(r.Context(), issues)
	r = r.WithContext(ctx)

	return false
}

func (s *Server) AddValidationErrors(w http.ResponseWriter, r *http.Request, errMap map[string]string) {
	issues := make(gonertia.ValidationErrors, len(errMap))
	for field, err := range errMap {
		issues[field] = err
	}

	ctx := gonertia.AddValidationErrors(r.Context(), issues)
	r = r.WithContext(ctx)
}

func (s *Server) AddInternalErrorMsg(w http.ResponseWriter, r *http.Request) {
	s.AddValidationErrors(w, r, map[string]string{
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

func (Server) RealIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		if i := strings.Index(ip, ","); i != -1 {
			ip = ip[:i]
		}

		return strings.TrimSpace(strings.Trim(ip, "[]"))
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return strings.Trim(ip, "[]")
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	return ip
}

func (s *Server) PrintRoutes() {
	chi.Walk(s.Mux, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		slog.Info("Route", slog.String("method", method), slog.String("route", route))
		return nil
	})
}

func NewServer(c Config, inertia *gonertia.Inertia) *Server {
	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
	)

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

	server := Server{
		Mux:            r,
		inertia:        inertia,
		cookieSecure:   c.Cookie.Secure,
		cookieHTTPOnly: c.Cookie.HTTPOnly,
		cookiePath:     c.Cookie.Path,
		cookieDomain:   c.Cookie.Domain,
		cookieSameSite: sameSite,
	}

	server.AddStaticAssetsRoute()

	return &server
}
