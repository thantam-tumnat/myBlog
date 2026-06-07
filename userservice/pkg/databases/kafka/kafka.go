package kafka

import (
	"userservice/configs"
	"userservice/modules/logs"
	"userservice/pkg/utils"

	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2/log"
)

func NewKafkaConsumerGroup(cfg *configs.Configs) (sarama.ConsumerGroup, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = cfg.Kafka.ClientID

	serversUrl1, _, err := utils.ConnectionUrlBuilder("kafka", cfg)
	if err != nil {
		return nil, err
	}
	consumer, err := sarama.NewConsumerGroup([]string{serversUrl1}, "tcchub.task", saramaConfig)
	if err != nil {
		return nil, err
	}
	logs.Info("kafka ConsumerGroup has been connected")
	return consumer, nil

}

func NewKafkaSyncProducer(cfg *configs.Configs) (sarama.SyncProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = cfg.Kafka.ClientID
	saramaConfig.Producer.Return.Successes = true
	serversUrl1, _, err := utils.ConnectionUrlBuilder("kafka", cfg)
	if err != nil {
		return nil, err
	}
	producer, err := sarama.NewSyncProducer([]string{serversUrl1}, saramaConfig)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	logs.Info("kafka SyncProducer has been connected")
	return producer, nil
}
