package main

import (
	"log"
	"strings"
	"blogservice/configs"
	"blogservice/modules/logs"
	"blogservice/modules/servers"
	
	kafka "blogservice/pkg/databases/kafka"
	// minio "task-service/pkg/databases/minio"
	databases "blogservice/pkg/databases/postgresql"
	redis "blogservice/pkg/databases/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	cfg := new(configs.Configs)
	logs.Info("Loading configuration...")
	configs.LoadConfigs(cfg)
	logs.Info("Configuration loaded successfully.")

	// New Database
	db, err := databases.NewPostgreSQLDBConnection(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
	//defer db.Close()

	// Kafka Producer Config
	producer, err := kafka.NewKafkaSyncProducer(cfg)
	if err != nil {
		logs.Error(err)
	}
	defer producer.Close()

	// Kafka Consumer Config
	consumer, err := kafka.NewKafkaConsumerGroup(cfg)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	redis := redis.InitRedis(cfg)
	// minio := minio.NewMinioClient(cfg)

	s := servers.NewServer(cfg, db, consumer, producer, redis)
	s.App.Use(logs.LogHTTPMethod)
	s.App.Use(cors.New(cors.Config{
		AllowOrigins: cfg.Cors.AllowOrigins,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowHeaders: "Content-Type, Authorization",
	}))
	s.Start()
}
