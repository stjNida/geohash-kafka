package main

import (
	"encoding/json"
	"fmt"
	"kafka-golang/geohash"
	"kafka-golang/producer"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

var geolocation producer.Geolocation

func main() {

	topic := "geolocations"
	worker, err := connectConsumer([]string{"localhost:9092"})
	if err != nil {
		panic(err)
	}

	// Calling ConsumePartition. It will open one connection per broker
	// and share it for all partitions that live on it.
	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	fmt.Println("Consumer started ")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// Count how many message processed
	msgCount := 0

	// Get signal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, string(msg.Topic), string(msg.Value))

				if topic == "geolocations" {
					json.Unmarshal(msg.Value, &geolocation)
					long := geolocation.Longitude
					lat := geolocation.Latitude
					fmt.Printf("Latitude: %f , Longitude: %f  \n", lat, long)
					geohash := createGeohash(lat, long)
					producer.PushToQueue("geohashes", []byte(geohash))

				}

			case <-sigchan:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")

	if err := worker.Close(); err != nil {
		panic(err)
	}

}

func connectConsumer(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create new consumer
	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func createGeohash(longitude, latitude float64) string {
	precision := 7
	geohash := geohash.Encode(longitude, latitude, precision)
	fmt.Printf("Geohash: %s \n", geohash)
	return geohash

}
