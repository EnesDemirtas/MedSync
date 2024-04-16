package userapi

import (
	"net/http"

	"github.com/EnesDemirtas/medisync/apis/services/warehouse/mid"
	"github.com/EnesDemirtas/medisync/app/api/authsrv"
	"github.com/EnesDemirtas/medisync/app/core/crud/userapp"
	"github.com/EnesDemirtas/medisync/business/api/auth"
	"github.com/EnesDemirtas/medisync/business/core/crud/userbus"
	"github.com/EnesDemirtas/medisync/foundation/logger"
	"github.com/EnesDemirtas/medisync/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	UserBus 	*userbus.Core
	AuthSrv    	*authsrv.AuthSrv
	Log   		*logger.Logger
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.Log, cfg.AuthSrv)
	ruleAdmin := mid.Authorize(cfg.Log, cfg.AuthSrv, auth.RuleAdminOnly)
	ruleAuthorizeUser := mid.AuthorizeUser(cfg.Log, cfg.AuthSrv, cfg.UserBus, auth.RuleAdminOrSubject)
	ruleAuthorizeAdmin := mid.AuthorizeUser(cfg.Log, cfg.AuthSrv, cfg.UserBus, auth.RuleAdminOnly)

	api := newAPI(userapp.NewCore(cfg.UserBus))
	app.Handle(http.MethodGet, version, "/users", api.query, authen, ruleAdmin)
	app.Handle(http.MethodGet, version, "/users/{user_id}", api.queryByID, authen, ruleAuthorizeUser)
	app.Handle(http.MethodPost, version, "/users", api.create, authen, ruleAdmin)
	app.Handle(http.MethodPut, version, "/users/role/{user_id}", api.updateRole, authen, ruleAuthorizeAdmin)
	app.Handle(http.MethodPut, version, "/users/{user_id}", api.update, authen, ruleAuthorizeUser)
	app.Handle(http.MethodDelete, version, "/users/{user_id}", api.delete, authen, ruleAuthorizeUser)
}