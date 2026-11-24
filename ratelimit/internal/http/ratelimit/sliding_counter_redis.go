package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"time"

	_ "embed"

	"github.com/redis/go-redis/v9"
)

var (
	//go:embed sliding_counter.lua
	slidingCounterScriptRaw string
	slidingCounterScript    = redis.NewScript(slidingCounterScriptRaw)
)

type RedisSlidingCounter struct {
	scripter redis.Scripter
	key      string
	window   time.Duration
	limit    int
}

func NewSlidingCounter(
	scripter redis.Scripter, key string, window time.Duration, limit int,
) *RedisSlidingCounter {
	return &RedisSlidingCounter{
		scripter: scripter,
		key:      key,
		window:   window,
		limit:    limit,
	}
}

func (s *RedisSlidingCounter) Allow(ctx context.Context, key string) (bool, error) {
	var (
		windowMs = s.window.Milliseconds()
		nowMs    = time.Now().UnixMilli()
	)

	curStartMs := nowMs - nowMs%windowMs
	prevStartMs := curStartMs - windowMs

	curKey := s.key + ":" + key + ":" + strconv.FormatInt(curStartMs, 10)
	prevKey := s.key + ":" + key + ":" + strconv.FormatInt(prevStartMs, 10)

	allowed, err := slidingCounterScript.Run(
		ctx,
		s.scripter,
		[]string{curKey, prevKey},
		windowMs,
		s.limit,
		nowMs,
	).Bool()
	if err != nil {
		return false, fmt.Errorf("run script: %w", err)
	}

	return allowed, nil
}
