package producer

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type Geolocation struct {
	Latitude  float64 `form:"latitude" json:"latitude"`
	Longitude float64 `form:"longitude" json:"longitude"`
}

// Connect to Apache Kafka using sarama as Producer
func ConnectProducer(brokersUrl []string) (sarama.SyncProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func PushToQueue(topic string, message []byte) error {

	brokersUrl := []string{"localhost:9092"}
	producer, err := ConnectProducer(brokersUrl)
	if err != nil {
		return err
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

	return nil
}
