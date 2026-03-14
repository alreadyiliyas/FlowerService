package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ilyas/flower/services/auth/internal/apperrors"
	authusecase "github.com/ilyas/flower/services/auth/internal/usecase/auth"
	"github.com/ilyas/flower/services/auth/internal/utils"
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

func AuthMiddleware(jwtSecret string, authUC authusecase.UsecaseAuth) func(http.Handler) http.Handler {
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

			claims, err := utils.ParseAccessToken(parts[1], jwtSecret)
			if err != nil || claims == nil || claims.SessionID == "" {
				utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
				return
			}

			if headerSessionID := r.Header.Get("X-Session-Id"); headerSessionID != "" && headerSessionID != claims.SessionID {
				utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
				return
			}

			session, err := authUC.GetSession(r.Context(), claims.SessionID)
			if err != nil {
				if errors.Is(err, apperrors.ErrUnauthorized) {
					utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
					return
				}
				utils.Send(w, http.StatusInternalServerError, nil, "internal server error")
				return
			}
			if session.UserID != claims.UserID {
				utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
				return
			}

			ctx := context.WithValue(r.Context(), ctxUserID, session.UserID)
			ctx = context.WithValue(ctx, ctxRole, session.Role)
			ctx = context.WithValue(ctx, ctxPhone, session.PhoneNumber)
			ctx = context.WithValue(ctx, ctxSessionID, session.SessionID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
