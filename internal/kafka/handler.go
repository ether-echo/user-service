package kafka

import (
	"encoding/json"

	"github.com/ether-echo/user-service/internal/domain"

	"github.com/IBM/sarama"
)

type MessageHandler interface {
	ProcessStart(user *domain.User) error
}

type ConsumerGroupHandler struct {
	handler MessageHandler
}

func NewConsumerGroupHandler(handler MessageHandler) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{handler: handler}
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Info("message channel was closed")
				return nil
			}

			var user domain.User

			log.Infof("Processing message: %s", message.Value)

			if err := json.Unmarshal(message.Value, &user); err != nil {
				log.Errorf("Error unmarshaling message: %v", err)
				continue
			}

			user.Command = message.Topic

			switch message.Topic {
			case "start":
				if err := h.handler.ProcessStart(&user); err != nil {
					log.Errorf("Error handling message: %v", err)
					continue
				}
			default:
				log.Infof("Unknown topic: %v", message.Topic)
			}

			session.MarkMessage(message, "processed in user service")

		case <-session.Context().Done():
			return nil
		}
	}
}
