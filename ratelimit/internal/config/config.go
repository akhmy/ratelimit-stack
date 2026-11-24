package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type RateLimitConfig struct {
	Key    string        `env:"RATE_LIMIT_KEY" envDefault:"rl"`
	Window time.Duration `env:"RATE_LIMIT_WINDOW" envDefault:"10s"`
	Limit  int           `env:"RATE_LIMIT_LIMIT" envDefault:"3"`
}

type Config struct {
	ServerAddr string `env:"SERVER_ADDR" envDefault:":8080"`
	RedisAddr  string `env:"REDIS_ADDR,required"`
	RL         RateLimitConfig
}

func ParseEnv() (Config, error) {
	return env.ParseAs[Config]()
}
