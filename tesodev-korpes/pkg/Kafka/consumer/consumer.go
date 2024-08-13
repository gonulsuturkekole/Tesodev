package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func RunConsumer(consumerAction func(map[string]interface{}), topic string, config map[string]string) { //Takes a function as a parameter, executes this function when message received
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        config["bootstrap.servers"],
		"sasl.mechanisms":          config["sasl.mechanisms"],
		"security.protocol":        config["security.protocol"],
		"sasl.username":            config["sasl.username"],
		"sasl.password":            config["sasl.password"],
		"group.id":                 config["group.id"],
		"auto.offset.reset":        config["auto.offset.reset"],
		"allow.auto.create.topics": config["allow.auto.create.topics"],
	})

	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics([]string{topic, "^aRegex.*[Tt]opic"}, nil)
	if err != nil {
		log.Print("CONSUMER ERROR: ")
		panic(err)
	}

	go func() {
		for {
			msg, err := c.ReadMessage(-1)
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
				hashmap := make(map[string]interface{})
				json.Unmarshal(msg.Value, &hashmap)
				consumerAction(hashmap) // RUN THE PASSED FUNCTION WHEN MESSAGE RECEIVED
			} else {
				// The client will automatically try to recover from all errors.
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}()
	// c.Close()
}
