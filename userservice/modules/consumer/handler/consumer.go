package consumer

import (
	"fmt"
	"userservice/modules/logs"

	"github.com/IBM/sarama"
)

type EventHandler interface {
	Handle(topic string, eventBytes []byte) error
}

type consumerHandler struct {
	eventHandler EventHandler
}

func NewConsumerHandler(eventHandler EventHandler) sarama.ConsumerGroupHandler {
	return &consumerHandler{eventHandler}
}

func (obj *consumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (obj *consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (obj *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		err := obj.eventHandler.Handle(message.Topic, message.Value)
		if err != nil {
			logs.Error(fmt.Sprintf("Error handling message: %v", err))
			return err
		}
		session.MarkMessage(message, "")
	}
	return nil
}
