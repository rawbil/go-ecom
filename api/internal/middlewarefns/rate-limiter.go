package middlewarefns

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/rawbil/ecom2/internal/utils"
	"golang.org/x/time/rate"
)

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	clients map[string]*clientLimiter
	mu      sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*clientLimiter),
	}
}

func (rl *RateLimiter) getClientLimiter(key string, r rate.Limit, burst int) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if client, exists := rl.clients[key]; exists {
		client.lastSeen = time.Now()
		return client.limiter
	}

	limiter := rate.NewLimiter(r, burst)
	rl.clients[key] = &clientLimiter{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	return limiter
}

func (rl *RateLimiter) Limit(ratePerTime rate.Limit, burst int) Mdlw {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := utils.GetClientIP(r)
			key := ip + ":" + r.URL.Path

			limiter := rl.getClientLimiter(key, ratePerTime, burst)

			allowed := limiter.Allow()

			utils.Log.Info(
				"Rate limit check",
				"ip", ip,
				"key", key,
				"path", r.URL.Path,
				"allowed", allowed,
			)

			if !allowed {
				utils.ErrorHandler(errors.New("Too many requests"), "Too many requests. Try again later", w, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (rl *RateLimiter) Cleanup(maxAge time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, client := range rl.clients {
		if now.Sub(client.lastSeen) > maxAge {
			delete(rl.clients, key)
		}
	}
}

func (rl *RateLimiter) StartCleanup(interval, maxAge time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		rl.Cleanup(maxAge)
	}
}
