package cmd

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	config2 "tesodev-korpes/CustomerService/config"
	"tesodev-korpes/CustomerService/internal"
	"tesodev-korpes/OrderService/client"
	"tesodev-korpes/pkg"
)

func BootCustomerService(client *mongo.Client, h_client *client.CustomerClient, e *echo.Echo) {
	config := config2.GetCustomerConfig("dev")
	customerCol, err := pkg.GetMongoCollection(client, config.DbConfig.DBName, config.DbConfig.ColName)
	if err != nil {
		panic(err)
	}
	addressCol, err := pkg.GetMongoCollection(client, config.DbConfig.DBName, "customeraddresses")
	if err != nil {
		panic(err)
	}
	phoneNumber, err := pkg.GetMongoCollection(client, config.DbConfig.DBName, "customerphone")
	if err != nil {
		panic(err)
	}
	repo := internal.NewRepository(customerCol, addressCol, phoneNumber)
	service := internal.NewService(repo, h_client)
	internal.NewHandler(e, service)

	e.Logger.Fatal(e.Start(config.Port))

}
