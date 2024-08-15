package cmd

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	config4 "tesodev-korpes/ConsumerService/config"
	"tesodev-korpes/ConsumerService/internal"
	"tesodev-korpes/pkg"
)

func BootConsumerService(client *mongo.Client, e *echo.Echo) {

	config := config4.GetConsumerConfig("dev")

	// MongoDB koleksiyonunu al
	consumerCol, err := pkg.GetMongoCollection(client, config.DbConfig.DBName, config.DbConfig.ColName)
	if err != nil {
		panic(err)
	}
	repo := internal.NewRepository(consumerCol)
	service := internal.NewService(repo)
	internal.NewHandler(e, service)
	e.Logger.Fatal(e.Start(config.Port))
}
