package servers

import (
	"blogservice/configs"
	"context"
	"fmt"

	"blogservice/modules/entities"
	"blogservice/modules/logs"
	"blogservice/pkg/utils"

	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type server struct {
	App                  *fiber.App
	Db                   *gorm.DB
	Cfg                  *configs.Configs
	ConsumerGroup        sarama.ConsumerGroup
	SyncProducer         sarama.SyncProducer
	consumerGroupHandler sarama.ConsumerGroupHandler
	Redis                *redis.Client
	// Minio                *minio.Client
}

func NewServer(cfg *configs.Configs,
	db *gorm.DB,
	consumerGroup sarama.ConsumerGroup,
	syncProducer sarama.SyncProducer,
	// Minio *minio.Client,
	redis *redis.Client,
) *server {
	return &server{
		App:           fiber.New(),
		Db:            db,
		Cfg:           cfg,
		ConsumerGroup: consumerGroup,
		SyncProducer:  syncProducer,
		// Minio:         Minio,
		Redis: redis,
	}
}

func (s *server) Start() {
	if err := s.Handlers(); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	fiberConnURL, _, err := utils.ConnectionUrlBuilder("fiber", s.Cfg)
	if err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	go func() {
		logs.Info("--- User Consumer Started ---")
		for {
			err := s.ConsumerGroup.Consume(context.Background(), entities.Topics, s.consumerGroupHandler)
			if err != nil {
				logs.Error(fmt.Sprintf("Kafka Consume ERROR : %v", err))
			}
		}
	}()

	port := s.Cfg.App.Port
	logs.Info(fmt.Sprintf("server started on port:%s", port))

	if err := s.App.Listen(fiberConnURL); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
}
