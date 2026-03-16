package middleware

import (
	"net/http"

	"github.com/ilyas/flower/services/catalog/internal/apperrors"
	"github.com/ilyas/flower/services/catalog/internal/utils"
)

func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := RoleFromContext(r.Context())
			if !ok || role == "" {
				utils.Send(w, http.StatusUnauthorized, nil, apperrors.ErrUnauthorized.Error())
				return
			}
			if _, exists := allowed[role]; !exists {
				utils.Send(w, http.StatusForbidden, nil, apperrors.ErrForbidden.Error())
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
