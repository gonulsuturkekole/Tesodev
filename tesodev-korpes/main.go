package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tesodev-korpes/ConsumerService/clientCon"
	consumerCmd "tesodev-korpes/ConsumerService/cmd"
	"tesodev-korpes/CustomerService/cmd"
	"tesodev-korpes/OrderService/client"
	_ "tesodev-korpes/OrderService/client"
	orderCmd "tesodev-korpes/OrderService/cmd"
	_ "tesodev-korpes/docs"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/Kafka/consumer"
	"tesodev-korpes/pkg/Kafka/producer"
	"tesodev-korpes/pkg/middlewares"
	"tesodev-korpes/shared/config"
)

type RequestProcessor struct {
	counter int
	mutex   sync.Mutex
}

// Increment increments the counter
func (rp *RequestProcessor) Increment() {
	rp.mutex.Lock()
	defer rp.mutex.Unlock()
	rp.counter++
}

// GetCounter returns the current value of the counter
func (rp *RequestProcessor) GetCounter() int {
	rp.mutex.Lock()
	defer rp.mutex.Unlock()
	return rp.counter
}

func main() {
	// Database configuration based on environment (dev, qa, prod)
	dbConf := config.GetDBConfig("dev")

	client1, err := pkg.GetMongoClient(dbConf.MongoDuration, dbConf.MongoClientURI)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)

	}
	h_client := client.NewCustomerClient(pkg.NewRestClient())
	consumerClient := clientCon.NewConsumerClient(pkg.NewRestClient())
	// Initialize Echo

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(pkg.CorrelationIDMiddleware)
	e.Use(middlewares.Logger())
	e.Use(pkg.Authenticate)

	// Kafka configuration settings
	brokers := []string{"localhost:9092"}
	topic := "your_topic_name"

	// Initialize Kafka Producer
	kafkaProducer := producer.NewProducer(brokers, topic)

	kafkaConsumer := &consumer.Consumer{}

	// Simulate a request by producing a message every 5 seconds
	/*go func() {
		for {
			message := "This is a test message"
			err := kafkaProducer.ProduceMessage(message)
			if err != nil {
				fmt.Printf("Error producing message: %v\n", err)
			} else {
				fmt.Println("Produced message to Kafka")
			}
			time.Sleep(5 * time.Second) // Adjust the frequency as needed
		}
	}()*/

	/*	// Define the consumer action to process messages
		consumerAction := func(msg string, err error) {
			if err != nil {
				fmt.Printf("Error consuming message: %v\n", err)
				return
			}
			fmt.Printf("Consumed message: %s\n", msg)
		}
	*/
	// Initialize Kafka Consumer
	/*kafkaConsumer := &consumer.Consumer{Topic: topic}
	kafkaConsumer.CreateConnection(brokers)*/

	// Handle shutdown signals (e.g., CTRL+C)
	/*	sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)*/

	// Handle shutdown signals (e.g., CTRL+C)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Determine which service(s) to start based on command-line argument
	if len(os.Args) < 2 {
		panic("Please provide a service to start: customer, order, or both")
	}
	input := os.Args[1] // This argument specifies which service to start.

	switch input {
	case "customer":
		cmd.BootCustomerService(client1, e)
	case "order":
		orderCmd.BootOrderService(client1, h_client, kafkaProducer, e)
	case "consumer":
		go consumerCmd.BootConsumerService(client1, kafkaConsumer, consumerClient, e, brokers, topic)
	case "both":
		go cmd.BootCustomerService(client1, e)
		go orderCmd.BootOrderService(client1, h_client, kafkaProducer, e)
		go consumerCmd.BootConsumerService(client1, kafkaConsumer, consumerClient, e, brokers, topic)
	default:
		panic("Invalid input. Use 'customer', 'order', or 'both'.")
	}

	<-sigs
	fmt.Println("Shutting down...")

	// Perform any necessary cleanup
	kafkaConsumer.Close()
	kafkaProducer.Close()

	fmt.Println("Kafka connections closed. Exiting.")
}
