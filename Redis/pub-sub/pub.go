package main

import (
	"context"
	"fmt"
	"os"

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

	// Publish messages to the 'chat' channel
	for {
		var message string
		fmt.Print("Enter message: ")
		fmt.Scanln(&message)
		if message == "exit" {
			break
		}
		err := rdb.Publish(ctx, "chat", message).Err()
		if err != nil {
			fmt.Println("Error publishing message:", err)
			os.Exit(1)
		}
	}

	// Close the Redis connection
	err := rdb.Close()
	if err != nil {
		fmt.Println("Error closing connection:", err)
		os.Exit(1)
	}
}
