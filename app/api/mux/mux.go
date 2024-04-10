// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"net/http"
	"os"

	midhttp "github.com/EnesDemirtas/medisync/app/api/mid/http"
	"github.com/EnesDemirtas/medisync/business/api/auth"
	"github.com/EnesDemirtas/medisync/business/core/crud/delegate"
	"github.com/EnesDemirtas/medisync/business/core/crud/inventorybus"
	"github.com/EnesDemirtas/medisync/business/core/crud/medicinebus"
	"github.com/EnesDemirtas/medisync/business/core/crud/tagbus"
	"github.com/EnesDemirtas/medisync/business/core/crud/userbus"
	"github.com/EnesDemirtas/medisync/foundation/logger"
	"github.com/EnesDemirtas/medisync/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

// Options represents optional parameters.
type Options struct {
	corsOrigin []string
}

// WithCORS provides configuration options for CORS.
func WithCORS(origins []string) func(opts *Options) {
	return func(opts *Options) {
		opts.corsOrigin = origins
	}
}

// BusCrud represents the set of core business packages.
type BusCrud struct {
	Delegate  *delegate.Delegate
	Inventory *inventorybus.Core
	Medicine  *medicinebus.Core
	Tag       *tagbus.Core
	User 	  *userbus.Core
}

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build 		string
	Shutdown    chan os.Signal
	Log 		*logger.Logger
	Auth 		*auth.Auth
	DB			*sqlx.DB
	Tracer		trace.Tracer
	BusCrud     BusCrud
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg Config)
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config, routeAdder RouteAdder, options ...func(opts *Options)) http.Handler {
	var opts Options
	for _, option := range options {
		option(&opts)
	}

	app := web.NewApp(
		cfg.Shutdown,
		cfg.Tracer,
		midhttp.Logger(cfg.Log),
		midhttp.Errors(cfg.Log),
		midhttp.Metrics(),
		midhttp.Panics(),
	)

	if len(opts.corsOrigin) > 0 {
		app.EnableCORS(midhttp.Cors(opts.corsOrigin))
	}

	routeAdder.Add(app, cfg)

	return app
}