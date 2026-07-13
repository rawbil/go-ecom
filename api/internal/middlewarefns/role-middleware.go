package middlewarefns

import (
	"net/http"
)


type Middleware func(f http.Handler) http.Handler

//! ROLE MIDDLEWARE
func RoleMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
		})
	}
}