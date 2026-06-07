package configs

import (
	"fmt"
	"os"
	"userservice/modules/logs"

	"github.com/joho/godotenv"
)

func LoadConfigs(cfg *Configs) {
	// Load dotenv config
	if err := godotenv.Load("../config.env"); err != nil {
		panic(err.Error())
	}
	// cfg = new(configs.Configs)

	//Fiber
	cfg.App.Host = os.Getenv("APP_HOST")
	cfg.App.Port = os.Getenv("USER_PORT")

	cfg.PostgreSQL.Host = os.Getenv("DB_HOST")
	cfg.PostgreSQL.Port = os.Getenv("DB_PORT")
	cfg.PostgreSQL.Username = os.Getenv("DB_USERNAME_USER")
	cfg.PostgreSQL.Password = os.Getenv("DB_PASSWORD_USER")
	cfg.PostgreSQL.Database = os.Getenv("DB_DATABASE_USER")
	cfg.PostgreSQL.Schema = os.Getenv("DB_SCHEMA_USER")
	cfg.PostgreSQL.SSLMode = os.Getenv("DB_SSLMOD")

	cfg.Kafka.ClientID = os.Getenv("KAFKA_CLIENT_ID")
	cfg.Kafka.Group = os.Getenv("KAFKA_GROUP_ID")
	cfg.Kafka.Host = os.Getenv("KAFKA_SERVERS")
	cfg.Kafka.Port = os.Getenv("KAFKA_PORT")

	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port = os.Getenv("REDIS_PORT")
	// cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

	logs.Info(fmt.Sprintf("Task App      --> Port : %v", cfg.App.Port))
	logs.Info(fmt.Sprintf("Task Database --> Port : %v , Driver : %v , Username : %v ,Schema : %v, SSLmod : %v", cfg.PostgreSQL.Port, cfg.PostgreSQL.Database, cfg.PostgreSQL.Username, cfg.PostgreSQL.Schema, cfg.PostgreSQL.SSLMode))
	logs.Info(fmt.Sprintf("Task Redis    --> Port : %v , Host : %v", cfg.Redis.Port, cfg.Redis.Host))
	logs.Info(fmt.Sprintf("Task Kafka    --> Port : %v , Group ID : %v , Client ID : %v", cfg.Kafka.Port, cfg.Kafka.Group, cfg.Kafka.ClientID))
}
