package main

import (
	"fmt"
	"io"
	"log"

	pb "github.com/mxngocqb/Golang/gRPC/greet/proto"
)

func (s *Server) LongGreet(strem pb.GreetService_LongGreetServer) error {
	log.Printf("LongGreet function was invoked")

	res := ""

	for {
		req, err := strem.Recv()

		if err == io.EOF {
			return strem.SendAndClose(&pb.GreetResponse{
				Result: res,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}

		res += fmt.Sprintf("Hello %s!\n", req.FirstName)
		
	}
}
