package authorization

import (
	"errors"
	"net/http"
	"slices"

	repository "github.com/rawbil/ecom2/internal/adapters/sqlc"
	authutils "github.com/rawbil/ecom2/internal/auth/auth-utils"
	"github.com/rawbil/ecom2/internal/utils"
)

type Middleware func(http.Handler) http.Handler

func PermissionMiddleware(repository repository.Queries, permission string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := authutils.GetUserFromContext(r.Context())
			if !ok {
				utils.ErrorHandler(errors.New("error getting context for role middleware"), "error getting context for role middleware", w, http.StatusUnauthorized)
				return
			}

			// // Get User
			// user, err := repository.ListUserById(r.Context(), claims.UserID)
			// if err != nil {
			// 	utils.ErrorHandler(err, "Failed to fetch user", w, http.StatusUnauthorized)
			// 	return
			// }

			// Get User Permissions
			permissions, err := repository.GetUserPermissions(r.Context(), claims.UserID)
			if err != nil {
				utils.ErrorHandler(err, "Error fetching permissions", w, http.StatusInternalServerError)
				return
			}

			if slices.Contains(permissions, permission) {
				next.ServeHTTP(w, r)
				return
			}

			utils.ErrorHandler(errors.New("Forbidden"), "Unauthorized resource", w, http.StatusForbidden)
		})
	}
}
