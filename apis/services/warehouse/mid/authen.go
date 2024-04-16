package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/EnesDemirtas/medisync/app/api/authsrv"
	"github.com/EnesDemirtas/medisync/app/api/errs"
	"github.com/EnesDemirtas/medisync/app/api/mid"
	"github.com/EnesDemirtas/medisync/foundation/logger"
	"github.com/EnesDemirtas/medisync/foundation/web"
)

// Authenticate validates a JWT from the `Authorization` header.
func Authenticate(log *logger.Logger, authsrv *authsrv.AuthSrv) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctxAuth, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()

			resp, err := authsrv.Authenticate(ctxAuth, r.Header.Get("authorization"))
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			ctx = mid.SetUserID(ctx, resp.UserID)
			ctx = mid.SetClaims(ctx, resp.Claims)

			return handler(ctx, w, r)
		}

		return h
	}

	return m
} 