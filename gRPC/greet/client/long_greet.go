package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Clement-Jean/grpc-go-course/greet/proto"
)

func doLongGreet(c pb.GreetServiceClient) {
	log.Println(("doLongGreet was invoked"))

	reqs := []*pb.GreetRequest{
		{FirstName: "Clement"},
		{FirstName: "Ngoc"},
		{FirstName: "Anh Duc"},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while calling LongGreet %v\n", err)
	}
	for _, req := range reqs {
		log.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while recevin response from LongGreet: %v\n", err)
	}

	log.Printf("LongGreet %s\n", res.Result)
}
