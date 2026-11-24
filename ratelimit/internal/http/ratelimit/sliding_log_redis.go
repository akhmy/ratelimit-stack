package ratelimit

import (
	"context"
	_ "embed"
	"fmt"
	"time"
	"uuid"

	"github.com/redis/go-redis/v9"
)

var (
	//go:embed sliding_log.lua
	slidingLogScriptRaw string
	slidingLogScript    = redis.NewScript(slidingLogScriptRaw)
)

type RedisSlidingLog struct {
	scripter redis.Scripter
	key      string
	window   time.Duration
	limit    int
}

func NewSlidingLog(scripter redis.Scripter, key string, window time.Duration, limit int) *RedisSlidingLog {
	return &RedisSlidingLog{
		scripter: scripter,
		key:      key,
		window:   window,
		limit:    limit,
	}
}

func (rl *RedisSlidingLog) Allow(ctx context.Context, key string) (bool, error) {
	ok, err := slidingLogScript.Eval(
		ctx,
		rl.scripter,
		[]string{rl.key + ":" + key},
		rl.window.Milliseconds(),
		rl.limit,
		time.Now().UnixMilli(),
		uuid.NewV4().String(),
	).Bool()
	if err != nil {
		return false, fmt.Errorf("eval script: %w", err)
	}

	return ok, nil
}
