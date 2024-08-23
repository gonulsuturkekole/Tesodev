package cmd

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"tesodev-korpes/ConsumerService/clientCon"
	config4 "tesodev-korpes/ConsumerService/config"
	"tesodev-korpes/ConsumerService/internal"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/Kafka/consumer"
)

func BootConsumerService(client *mongo.Client, kafkaConsumer *consumer.Consumer, conClient *clientCon.ConsumerClient, e *echo.Echo, brokers []string, topic string) {

	config := config4.GetConsumerConfig("dev")
	consumerCol, err := pkg.GetMongoCollection(client, config.DbConfig.DBName, config.DbConfig.ColName)
	if err != nil {
		panic(err)
	}
	repo := internal.NewFinanceRepository(consumerCol)
	service := internal.NewService(repo, conClient, kafkaConsumer, brokers, topic)
	key := config.DbConfig.SecretKey
	ctx := context.Background()
	consumerAction := func(msg string, err error) {
		if err != nil {
			fmt.Printf("Error consuming message: %v\n", err)
			return
		}
		fmt.Printf("Consumed message: %s\n", msg)

		err = service.ProcessMessage(ctx, msg, key)
		if err != nil {
			fmt.Printf("Error processing message: %v\n", err)
		}
	}

	go kafkaConsumer.Read(consumerAction)

	e.Logger.Fatal(e.Start(config.Port))
}
