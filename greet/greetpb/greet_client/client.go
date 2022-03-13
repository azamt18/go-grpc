package main

import (
	"context"
	"fmt"
	"github.com/simplesteph/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	fmt.Printf("hello from client\n")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("error in closing connection: %v", err)
		}
	}(conn)

	c := greetpb.NewGreetServiceClient(conn)
	//fmt.Printf("created client: %f", c)

	makeUnaryCall(c)
}

func makeUnaryCall(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do a unary rpc\n")

	// prepare request
	//c.Greet(ctx context.Context, in *GreetRequest, opts ...grpc.CallOption) (*GreetResponse, error)
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Johnson",
		},
	}

	// make a request
	res, err := c.Greet(context.Background(), req)

	// handle an error in response
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	// handle res
	log.Printf("Response from Greet: %v\n", res.Result)
}
