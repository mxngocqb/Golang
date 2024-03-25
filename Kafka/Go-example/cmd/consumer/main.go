package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/Shopify/sarama"
)

func main() {
	brokers := []string{"kafka:9092", "kafka2:9093"} // Danh sách các broker

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Tạo consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
	defer consumer.Close()

	// Danh sách các topic để subscribe
	topics := []string{"topic1", "topic2"}

	// Chạy consumer cho mỗi topic
	var wg sync.WaitGroup
	for _, topic := range topics {
		wg.Add(1)
		go func(topic string) {
			defer wg.Done()

			partitions, err := consumer.Partitions(topic)
			if err != nil {
				log.Printf("Failed to get partitions for topic %s: %v", topic, err)
				return
			}

			for _, partition := range partitions {
				partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
				if err != nil {
					log.Printf("Failed to start partition consumer for topic %s, partition %d: %v", topic, partition, err)
					continue
				}

				defer partitionConsumer.Close()

				// Đọc messages từ partition
				go func(pc sarama.PartitionConsumer) {
					for msg := range pc.Messages() {
						log.Printf("Received message from topic %s, partition %d: %s", topic, partition, string(msg.Value))
					}
				}(partitionConsumer)
			}
		}(topic)
	}

	// Đợi tín hiệu để kết thúc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	select {
	case <-sigchan:
		log.Println("Interrupt signal received, shutting down...")
		wg.Wait()
		log.Println("Shutdown complete.")
	}
}
