package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

func main() {
	brokers := []string{"kafka:9092", "kafka2:9093"} // Danh sách các broker

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// Tạo producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}
	defer producer.Close()

	// Danh sách các topic để gửi tin nhắn
	topics := []string{"topic1", "topic2"}

	// Gửi tin nhắn đến từng topic
	for _, topic := range topics {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder("Hello from producer!"),
		}

		// Gửi tin nhắn
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message to topic %s: %v", topic, err)
		} else {
			log.Printf("Message sent to topic %s, partition %d, offset %d", topic, partition, offset)
		}
	}

	// Đợi tín hiệu để kết thúc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	<-sigchan
	log.Println("Interrupt signal received, shutting down...")
}
