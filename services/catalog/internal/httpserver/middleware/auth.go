package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
	authclient "github.com/ilyas/flower/services/catalog/internal/grpc/authclient"
	"github.com/ilyas/flower/services/catalog/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	ctxUserID    contextKey = "user_id"
	ctxRole      contextKey = "role"
	ctxPhone     contextKey = "phone"
	ctxSessionID contextKey = "session_id"
)

func UserIDFromContext(ctx context.Context) (uint64, bool) {
	v, ok := ctx.Value(ctxUserID).(uint64)
	return v, ok
}

func RoleFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxRole).(string)
	return v, ok
}

func PhoneFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxPhone).(string)
	return v, ok
}

func SessionIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxSessionID).(string)
	return v, ok
}

func AuthMiddleware(authClient authclient.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
				return
			}

			resp, err := authClient.GetUserContext(r.Context(), parts[1], r.Header.Get("X-Session-Id"))
			if err != nil {
				log.Printf("| middleware | failed to get user: %v", err)
				switch status.Code(err) {
				case codes.InvalidArgument:
					utils.Send(w, http.StatusBadRequest, nil, err.Error())
				case codes.Unauthenticated:
					utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
				default:
					utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
				}
				return
			}

			ctx := context.WithValue(r.Context(), ctxUserID, resp.UserID)
			ctx = context.WithValue(ctx, ctxRole, resp.Role)
			ctx = context.WithValue(ctx, ctxPhone, resp.PhoneNumber)
			ctx = context.WithValue(ctx, ctxSessionID, resp.SessionID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
