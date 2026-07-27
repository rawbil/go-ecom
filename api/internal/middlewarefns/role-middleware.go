package middlewarefns

import (
	"errors"
	"net/http"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
	"github.com/rawbil/ecom2/internal/utils"
)

type Middleware func(http.Handler) http.Handler

func RoleMiddleware(repository repository.Queries, roles ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//& Get role from context
			claims, ok := authutils.GetUserFromContext(r.Context())
			if !ok {
				utils.ErrorHandler(errors.New("error getting context for role middleware"), "error getting context for role middleware", w, http.StatusUnauthorized)
				return
			}

			user, err := repository.ListUserById(r.Context(), claims.UserID)
			if err != nil {
				utils.ErrorHandler(err, "Failed to fetch user", w, http.StatusUnauthorized)
				return
			}

			user_role := user.Role

			allowed_roles := make(map[string]bool)

			for _, role := range roles {
				allowed_roles[role] = true
			}

			if !allowed_roles[user_role] {
				utils.ErrorHandler(errors.New("unauthorized resource"), "unauthorized resource", w, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}
