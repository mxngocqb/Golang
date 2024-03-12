package main

import (
	"log"

	pb "github.com/Clement-Jean/grpc-go-course/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr string = "0.0.0.0:50051"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}

	log.Print(conn.Target())
	
	defer conn.Close()
	c := pb.NewGreetServiceClient(conn)

	// doGreet(c)
	// doGreetManyTime(c)
	// doLongGreet(c)

	doGreetEveryon(c)
}
