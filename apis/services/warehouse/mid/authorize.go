package mid

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/EnesDemirtas/medisync/app/api/authsrv"
	"github.com/EnesDemirtas/medisync/app/api/errs"
	"github.com/EnesDemirtas/medisync/app/api/mid"
	"github.com/EnesDemirtas/medisync/business/domain/userbus"
	"github.com/EnesDemirtas/medisync/foundation/logger"
	"github.com/EnesDemirtas/medisync/foundation/web"
	"github.com/google/uuid"
)

// ErrInvalidID represents a condition where the id is not a uuid.
var ErrInvalidID = errors.New("ID is not in its proper form")

// Authorize executes the specified role and does not extract any domain data.
func Authorize(log *logger.Logger, authSrv *authsrv.AuthSrv, rule string) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			userID, err := mid.GetUserID(ctx)
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			ctxAuth, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()

			auth := authsrv.Authorize{
				Claims: mid.GetClaims(ctx),
				UserID: userID,
				Rule:   rule,
			}

			if err := authSrv.Authorize(ctxAuth, auth); err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// AuthorizeUser executes the specified role and extracts the specified user
// from the DB if a user id is specified in the call. Depending on the rule
// specified, the userid from the claims may be compared with the specified
// user id.
func AuthorizeUser(log *logger.Logger, authSrv *authsrv.AuthSrv, userBus *userbus.Core, rule string) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var userID uuid.UUID

			if id := web.Param(r, "user_id"); id != "" {
				var err error
				userID, err = uuid.Parse(id)
				if err != nil {
					return errs.New(errs.Unauthenticated, ErrInvalidID)
				}

				usr, err := userBus.QueryByID(ctx, userID)
				if err != nil {
					switch {
					case errors.Is(err, userbus.ErrNotFound):
						return errs.New(errs.Unauthenticated, err)
					default:
						return errs.Newf(errs.Unauthenticated, "querybyid: userID[%s]: %s", userID, err)
					}
				}

				ctx = mid.SetUser(ctx, usr)
			}

			ctxAuth, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()

			auth := authsrv.Authorize{
				Claims: mid.GetClaims(ctx),
				UserID: userID,
				Rule:   rule,
			}

			if err := authSrv.Authorize(ctxAuth, auth); err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}