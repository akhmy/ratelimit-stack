package main

import (
	"log/slog"
	"net/http"
	"os"
	"ratelimit/internal/config"
	inhttp "ratelimit/internal/http"

	"github.com/redis/go-redis/v9"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.ParseEnv()
	if err != nil {
		logger.Error("failed to parse config", "error", err)
		return 1
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: inhttp.NewRouter(rdb, cfg.RL.Key, cfg.RL.Window, cfg.RL.Limit),
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server", "error", err)
		return 1
	}

	return 0
}
