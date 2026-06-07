package redis

import (
	"blogservice/configs"
	"blogservice/modules/logs"
	"blogservice/pkg/utils"

	"github.com/go-redis/redis/v8"
)

func InitRedis(cfg *configs.Configs) *redis.Client {
	url, _, err := utils.ConnectionUrlBuilder("redis", cfg)
	if err != nil {
		logs.Error(err)
	}
	return redis.NewClient(&redis.Options{
		Addr:     url,
		// Password: cfg.Redis.Password,
	})
}
