package utils

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

var clients = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := clients[ip]; !exists {
		clients[ip] = rate.NewLimiter(2, 5)
	}
	return clients[ip]
}

func RateLimiter(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetClientIP(r)
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded. Try again later.", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
