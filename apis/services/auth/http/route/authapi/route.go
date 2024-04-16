package authapi

import (
	"net/http"

	"github.com/EnesDemirtas/medisync/apis/services/auth/http/mid"
	"github.com/EnesDemirtas/medisync/app/core/crud/userapp"
	"github.com/EnesDemirtas/medisync/business/api/auth"
	"github.com/EnesDemirtas/medisync/business/core/crud/userbus"
	"github.com/EnesDemirtas/medisync/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	UserBus *userbus.Core
	Auth    *auth.Auth
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.UserBus, cfg.Auth)

	api := newAPI(userapp.NewCoreWithAuth(cfg.UserBus, cfg.Auth), cfg.Auth)
	app.Handle(http.MethodGet, version, "/auth/token/{kid}", api.token, authen)
	app.Handle(http.MethodGet, version, "/auth/authenticate", api.authenticate, authen)
	app.Handle(http.MethodPost, version, "/auth/authorize", api.authorize)
}