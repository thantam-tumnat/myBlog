package producer

import (
	"encoding/json"
	"fmt"
	"userservice/modules/entities"
	"userservice/modules/logs"

	// "github.com/Shopify/sarama"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type eventProducer struct {
	producer sarama.SyncProducer
}

func NewEventProducer(producer sarama.SyncProducer) entities.EventProducer {
	return eventProducer{producer}
}

func (obj eventProducer) Produce(event entities.Event) error {
	topic := event.TopicName()


	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	_, _, err = obj.producer.SendMessage(&msg)
	if err != nil {
		return err
	}
	logs.Info("producer message",
		zap.String("message", fmt.Sprintf("%s (success)", topic)),
	)

	return nil
}

