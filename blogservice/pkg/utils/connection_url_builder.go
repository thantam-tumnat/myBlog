package utils

import (
	"errors"
	"fmt"
	"blogservice/configs"
)

func ConnectionUrlBuilder(stuff string, cfg *configs.Configs) (string, string, error) {
	var url string
	var url2 string

	switch stuff {
	case "fiber":
		url = fmt.Sprintf(":%s", cfg.App.Port)
	case "postgresql":
		url = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.PostgreSQL.Host,
			cfg.PostgreSQL.Port,
			cfg.PostgreSQL.Username,
			cfg.PostgreSQL.Password,
			cfg.PostgreSQL.Database,
			cfg.PostgreSQL.SSLMode,
		)
		url2 = fmt.Sprintf("set search_path=%v", cfg.PostgreSQL.Schema)
	case "redis":
		url = fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	// case "minio":
	// 	url = fmt.Sprintf("%s:%s", cfg.Minio.Host, cfg.Minio.Port)
	case "kafka":
		url = fmt.Sprintf("%s:%s", cfg.Kafka.Host, cfg.Kafka.Port)
	default:
		errMsg := fmt.Sprintf("error, connection url builder doesn't know the %s", stuff)
		return "", "", errors.New(errMsg)
	}
	return url, url2, nil
}
