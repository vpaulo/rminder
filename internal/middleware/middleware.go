package middleware

import (
	"log"
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

// MiddlewareStack chains multiple middlewares
func MiddlewareStack(ms ...Middleware) Middleware {
	return Middleware(func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			m := ms[i]
			next = m(next)
		}

		return next
	})
}

func IsLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("checking if is logged in...")

		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)
		log.Println(time.Since(start), r.Method, r.URL.Path)
	})
}
