package cmd

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"tesodev-korpes/OrderService/client"
	config3 "tesodev-korpes/OrderService/config"
	"tesodev-korpes/OrderService/internal"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/Kafka/producer"
)

func BootOrderService(client *mongo.Client, h_client *client.CustomerClient, kafkaProducer *producer.Producer, e *echo.Echo) {
	config := config3.GetOrderConfig("dev")
	orderCol, err := pkg.GetMongoCollection(client, config.DbConfig.DBName, config.DbConfig.ColName)
	if err != nil {
		panic(err)
	}

	repo := internal.NewRepository(orderCol)
	service := internal.NewService(repo, h_client, kafkaProducer)
	internal.NewHandler(e, service)

	e.Logger.Fatal(e.Start(config.Port))
}
