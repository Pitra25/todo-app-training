package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ConfigRedis struct {
	Addr        string        `yml:"addr"`
	Password    string        `yml:"password"`
	User        string        `yml:"user"`
	DB          int           `yml:"db"`
	MaxRetries  int           `yml:"max_retries"`
	DialTimeout time.Duration `yml:"dial_timeout"`
	Timeout     time.Duration `yml:"timeout"`
}

func New(cfg ConfigRedis) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	ctx := context.Background()
	if err := db.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return db, nil
}
