// Package all binds all the routes into the specified app.
package all

import (
	"github.com/EnesDemirtas/medisync/apis/services/auth/mux"
	"github.com/EnesDemirtas/medisync/apis/services/auth/route/authapi"
	"github.com/EnesDemirtas/medisync/apis/services/auth/route/checkapi"
	"github.com/EnesDemirtas/medisync/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// the RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouteAdder interface.
func (add) Add(app *web.App, cfg mux.Config) {
	checkapi.Routes(app, checkapi.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:	   cfg.DB,
	})

	authapi.Routes(app, authapi.Config{
		UserBus: cfg.BusDomain.User,
		Auth:	 cfg.Auth,
	})
}