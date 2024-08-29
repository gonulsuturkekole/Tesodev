package cmd

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"tesodev-korpes/OrderService/client"
	config3 "tesodev-korpes/OrderService/config"
	"tesodev-korpes/OrderService/internal"
	_ "tesodev-korpes/docs"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/Kafka/producer"
)

// @title Customer Service API
// @version 1.0
// @description API documentation for Customer Service.
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8002
// @BasePath /api/v1
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
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
