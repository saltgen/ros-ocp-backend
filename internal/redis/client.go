package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"

	"github.com/redhatinsights/ros-ocp-backend/internal/config"
)

var (
	client *redis.Client
	cfg    *config.Config = config.GetConfig()
	once   sync.Once
)

func Init() error {
	var initErr error
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
			Username: cfg.RedisUsername,
			Password: cfg.RedisPassword,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if _, err := client.Ping(ctx).Result(); err != nil {
			initErr = fmt.Errorf("redis connection failed: %w", err)
		}
	})
	return initErr
}

func Client() *redis.Client {
	if err := Init(); err != nil {
		log.Error(err)
	}
	return client
}
