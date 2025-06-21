package rabbitmq

import (
	"encoding/json"

	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"

	"notification-service/config"
	"notification-service/internal/core/domain/entity"
)

type ConsumeRabbitmqInterface interface {
	ConsumeMessage(queueName string, req entity.NotificationEntity) error
}

type consumeRabbitMQ struct {
	conn *amqp.Connection
}

// ConsumeMessage implements ConsumeRabbitmqInterface.
func (c consumeRabbitMQ) ConsumeMessage(queueName string, req entity.NotificationEntity) error {
	ch, err := c.conn.Channel()
	if err != nil {
		log.Errorf("[ConsumeMessage-1] failed to create channel: %v", err)
		return err
	}

	defer ch.Close()

	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Errorf("[ConsumeMessage-2] failed to consume message: %v", err)
		return err
	}

	for msg := range msgs {
		log.Infof("[ConsumeMessage] received message: %s", msg.Body)
		if err = json.Marshal(msg.Body, &req); err != nil {
			log.Errorf("[ConsumeMessage-3] failed to unmarshal message: %v", err)
			continue
		}
	}

	return nil

}

func NewConsumeRabbitMQ(cfg config.Config) ConsumeRabbitmqInterface {

	newConnect, err := cfg.NewRabbitMQ()
	if err != nil {
		log.Fatalf("failed to create ConsumeRabbitMQ")
		return nil
	}

	defer newConnect.Close()

	return &consumeRabbitMQ{
		conn: newConnect,
	}
}
