package middlewarefns

import (
	"errors"
	"net/http"

	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
	"github.com/rawbil/ecom2/internal/utils"
)

type Middleware func(http.Handler) http.Handler

func RoleMiddleware(roles ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//& Get role from context
			claims, ok := authutils.GetUserFromContext(r.Context())
			if !ok {
				utils.Log.Error("error getting context for role middleware")
				return
			}

			user_role := claims.UserRole

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
