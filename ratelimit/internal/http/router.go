package http

import (
	"net/http"
	"ratelimit/internal/http/ratelimit"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRouter(rdb redis.Scripter, key string, window time.Duration, limit int) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /sliding-log", RateLimit(
		ratelimit.NewSlidingLog(rdb, key+":sl", window, limit),
	)(http.HandlerFunc(handler)))

	mux.Handle("GET /sliding-counter", RateLimit(
		ratelimit.NewSlidingCounter(rdb, key+":sc", window, limit),
	)(http.HandlerFunc(handler)))

	return mux
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
