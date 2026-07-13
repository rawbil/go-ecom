package middlewarefns

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/rawbil/ecom2/internal/utils"
)

type mw func(http.Handler) http.Handler

// ! Logger
func Logger() mw {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			next.ServeHTTP(w, r)

			utils.Log.Info(
				"HTTP Response",
				slog.Duration("duration", time.Since(start)),
				"method", r.Method,
				"path", r.URL.Path,
				"ip", r.RemoteAddr,
			)

		})
	}
}
