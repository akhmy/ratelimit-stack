package http

import (
	"context"
	"net/http"
)

type rateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
}

var tooManyRequestsText = http.StatusText(http.StatusTooManyRequests)

func RateLimit(rl rateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ok, err := rl.Allow(r.Context(), "42")
			if err != nil {
				panic(err)
			}
			if !ok {
				http.Error(w, tooManyRequestsText, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
