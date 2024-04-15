package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/EnesDemirtas/medisync/apis/services/warehouse/build/crud"
	"github.com/EnesDemirtas/medisync/app/api/mux"
	"github.com/EnesDemirtas/medisync/business/api/auth"
	"github.com/EnesDemirtas/medisync/business/core/crud/delegate"
	"github.com/EnesDemirtas/medisync/business/core/crud/userbus"
	"github.com/EnesDemirtas/medisync/business/core/crud/userbus/stores/userdb"
	"github.com/EnesDemirtas/medisync/business/data/sqldb"
	"github.com/EnesDemirtas/medisync/foundation/keystore"
	"github.com/EnesDemirtas/medisync/foundation/logger"
	"github.com/EnesDemirtas/medisync/foundation/web"
)

var build = "develop"

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******** SEND ALERT ********")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "WAREHOUSE-API", traceIDFn, events)

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	cfg := struct {
		// conf.Version
		Web struct {
			ReadTimeout 		time.Duration `conf:"default:5s"`
			WriteTimeout		time.Duration `conf:"default:10s"`
			IdleTimeout			time.Duration `conf:"default:120s"`
			ShutdownTimeout		time.Duration `conf:"default:20s"`
			APIHost 			string		  `conf:"default:0.0.0.0:3000"`
			DebugHost			string		  `conf:"default:0.0.0.0:4000"`
			CORSAllowedOrigins  []string	  `conf:"default:*"`
		}
		Auth struct {
			KeysFolder string `conf:"default:zarf/keys/"`
			ActiveKID  string `conf:"default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"`
			Issuer     string `conf:"default:service project"`
		}
		DB struct {
			User 		 string `conf:"default:postgres"`
			Password	 string `conf:"default:postgres,mask"`
			HostPort	 string `conf:"default:0.0.0.0:5432"`
			Name		 string `conf:"default:postgres"`
			MaxIdleConns int	`conf:"default:2"`
			MaxOpenConns int 	`conf:"default:0"`
			DisableTLS	 bool	`conf:"default:true"`
		}		
	}{}

	log.Info(ctx, "starting service")
	defer log.Info(ctx, "shutdown complete")

	log.Info(ctx, "startup", "status", "initializing database support", "hostport", cfg.DB.HostPort)

	db, err := sqldb.Open(sqldb.Config{
		User: 		  cfg.DB.User,
		Password:     cfg.DB.Password,
		HostPort:     cfg.DB.HostPort,
		Name:		  cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	defer db.Close()

	log.Info(ctx, "startup", "status", "initializing authentication support")

	ks := keystore.New()
	if err := ks.LoadRSAKeys(os.DirFS(cfg.Auth.KeysFolder)); err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}

	authCfg := auth.Config {
		Log: 		log,
		DB:			db,
		KeyLookup:  ks,
	}

	auth, err := auth.New(authCfg)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	log.Info(ctx, "startup", "status", "initializing business support")

	delegate := delegate.New(log)
	userBus := userbus.NewCore(log, delegate, userdb.NewStore(log, db))


	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := mux.Config{
		Build: build,
		Shutdown: shutdown,
		Log: log,
		Auth: auth,
		DB: db,
		BusCrud: mux.BusCrud{
			Delegate: delegate,
			User: userBus,
		},
	}

	api := http.Server{
		Addr: cfg.Web.APIHost,
		Handler: mux.WebAPI(cfgMux, buildRoutes(), mux.WithCORS(cfg.Web.CORSAllowedOrigins)),
		ReadTimeout: cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout: cfg.Web.IdleTimeout,
		ErrorLog: logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <- serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <- shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func buildRoutes() mux.RouteAdder {
	return crud.Routes()
}

