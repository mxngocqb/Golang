package main

import (
	"log"

	pb "github.com/mxngocqb/Golang/gRPC/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var addr string = "localhost:50051"

func main() {
	tls := true // change that to true if needed
	opts := []grpc.DialOption{}

	if tls {
		certFile := "ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certFile, "")

		if err != nil {
			log.Fatalf("Error while loading CA trust certificate: %v\n", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		creds := grpc.WithTransportCredentials(insecure.NewCredentials())
		opts = append(opts, creds)
	}

	conn, err := grpc.Dial(addr, opts...)

	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}

	log.Print(conn.Target())

	defer conn.Close()
	c := pb.NewGreetServiceClient(conn)

	doGreet(c)
	// doGreetManyTime(c)
	// doLongGreet(c)
	// doGreetEveryon(c)
	// doGreetWithDeadline(c, 1)
}
