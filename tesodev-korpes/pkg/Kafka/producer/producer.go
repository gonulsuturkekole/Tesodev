package producer

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func initiateProducer(config map[string]string) *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config["bootstrap.servers"],
		"sasl.mechanisms":   config["sasl.mechanisms"],
		"security.protocol": config["security.protocol"],
		"sasl.username":     config["sasl.username"],
		"sasl.password":     config["sasl.password"],
	})
	if err != nil {
		panic(err)
	}
	// defer p.Close()
	return p

}

type KafkaHandler struct {
	p *kafka.Producer
}

func (kf KafkaHandler) ProduceMessage(topic string, message map[string]interface{}) { // topic := "myTopic"
	// topic := "binance"
	hashmap, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	kf.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(hashmap),
	}, nil)
	fmt.Print("Message produced")
	// Wait for message deliveries before shutting down
	// kf.p.Flush(15 * 1000)

}

func NewKafkaHandler(config map[string]string) *KafkaHandler {
	kfHandler := initiateProducer(config)

	return &KafkaHandler{
		p: kfHandler,
	}
}
