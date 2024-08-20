package cmd

/*
import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"tesodev-korpes/pkg"
)


import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"tesodev-korpes/ConsumerService/clientCon"
	config4 "tesodev-korpes/ConsumerService/config"
	"tesodev-korpes/ConsumerService/internal"
	"tesodev-korpes/pkg"
)

func BootConsumerService(client *mongo.Client, clientCon *clientCon.ConsumerClient, e *echo.Echo) {

	config := config4.GetConsumerConfig("dev")

	consumerCol, err := pkg.GetMongoCollection(client, config.DbConfig.DBName, config.DbConfig.ColName)
	if err != nil {
		panic(err)
	}

	repo := internal.NewRepository(consumerCol)
	service := internal.NewService(repo, clientCon)
	internal.NewHandler(e, service)
	e.Logger.Fatal(e.Start(config.Port))
}
*/
