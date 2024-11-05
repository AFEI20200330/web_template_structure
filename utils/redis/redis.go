package redis

import (
	"fmt"
	"web_template/settings"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var rdb *redis.Client

// Init redis client
func Init(cfg *settings.RedisConfig) (err error) {
	// redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	// ping redis
	_, err = rdb.Ping().Result()
	if err != nil {
		zap.L().Error("redis connect error", zap.Error(err))
		return err
	}
	return nil
}

func Close() {
	if err := rdb.Close(); err != nil {
		zap.L().Error("redis close error", zap.Error(err))
	}
}
