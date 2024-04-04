package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/EnesDemirtas/medisync/app/api/errs"
	"github.com/EnesDemirtas/medisync/app/api/mid"
	"github.com/EnesDemirtas/medisync/business/api/auth"
	"github.com/EnesDemirtas/medisync/business/core/crud/userbus"
	"github.com/EnesDemirtas/medisync/foundation/web"
	"github.com/google/uuid"
)

// AuthorizeUser executes the specified rule and extracts the specified user
// from the DB if a user id is specified in the call. Depending on the rule
// specified, the userid from the claims may be compared with the specified
// user id.
func AuthorizeUser(a *auth.Auth, userBus *userbus.Core, rule string) web.MidHandler {
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

			claims := mid.GetClaims(ctx)
			if err := a.Authorize(ctx, claims, userID, rule); err != nil {
				return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
