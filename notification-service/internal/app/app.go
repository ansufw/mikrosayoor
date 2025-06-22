package app

import (
	"github.com/labstack/echo/v4"

	"notification-service/config"
	"notification-service/internal/adapter/message"
	"notification-service/internal/adapter/rabbitmq"
	"notification-service/utils"
)

func RunServer() {
	cfg := config.NewConfig()
	emailMessage := message.NewMessageEmail(cfg)
	rabbitMQAdapter := rabbitmq.NewConsumeRabbitMQ(cfg, emailMessage)

	e := echo.New()

	go func() {
		err := rabbitMQAdapter.ConsumeMessage(utils.NOTIF_EMAIL_VERIFICATION)
		if err != nil {
			e.Logger.Fatalf("failed to consume RabbitMQ for %s: %v", utils.NOTIF_EMAIL_VERIFICATION, err)
		}
	}()

	e.Logger.Fatal(e.Start(":" + cfg.App.AppPort))

}
