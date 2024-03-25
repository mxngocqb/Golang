package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()

	// Subscribe to the 'chat' channel
	pubsub := rdb.Subscribe(ctx, "chat")
	defer pubsub.Close()

	// Channel to listen for OS signals
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Listen for incoming messages
	ch := pubsub.Channel()
	for {
		select {
		case <-sigchan:
			fmt.Println("Received termination signal. Exiting...")
			return
		case msg := <-ch:
			fmt.Printf("Received: %s\n", msg.Payload)
			if msg.Payload == "exit" {
				fmt.Println("Received exit message. Exiting...")
				return
			}
		}
	}
}
