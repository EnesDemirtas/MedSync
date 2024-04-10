// Package crud binds the crud domain set of routes into the specified app.
package crud

import (
	"github.com/EnesDemirtas/medisync/apis/services/warehouseapiweb/routes/crud/userapi"
	"github.com/EnesDemirtas/medisync/app/api/mux"
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
	userapi.Routes(app, userapi.Config{
		UserBus: cfg.BusCrud.User,
		Auth:	 cfg.Auth,
	})
}