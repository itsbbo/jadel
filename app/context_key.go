package app

import (
	"net/http"

	"github.com/itsbbo/jadel/gonertia"
	"github.com/itsbbo/jadel/model"
)

const (
	CSRFCookieName = "XSRF-TOKEN"
	CSRFHeaderName = "X-XSRF-TOKEN"
	SessionKey     = "jadel_session"
	UserKey        = "user"
)

func CurrentUser(r *http.Request) *model.User {
	props := gonertia.PropsFromContext(r.Context())
	return props[UserKey].(*model.User)
}
