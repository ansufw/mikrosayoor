package message

import (
	"encoding/json"

	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"

	"user-service/config"
)

func PublishMessage(email, msg, notifType string) error {
	conn, err := config.NewConfig().NewRabbitMQ()
	if err != nil {
		log.Errorf("[PublishMessage-1] failed to connect to RabbitMQ: %v", err)
		return err
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[PublishMessage-2] failed to open a channel: %v", err)
		return err
	}

	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"notifications",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Errorf("[PublishMessage-3] failed to declare a queue: %v", err)
		return err
	}

	notification := map[string]string{
		"email": email,
		"msg":   msg,
	}

	body, err := json.Marshal(notification)
	if err != nil {
		log.Errorf("[PublishMessage-4] failed to marshal notification: %v", err)
		return err
	}

	return ch.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
