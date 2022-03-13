package main

import (
	"context"
	"fmt"
	"github.com/simplesteph/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"io"
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
	makeServerStreamingCall(c)
}

func makeUnaryCall(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do a unary rpc\n")

	// prepare unary request
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

	log.Printf("\t End of unary call \n\n")
}

func makeServerStreamingCall(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting server streaming RPC\n")

	// prepare streaming request
	request := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Johnson",
		},
	}

	// make a stream request
	responseStream, error := c.GreetManyTimes(context.Background(), request)

	// handle an error
	if error != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", error)
	}

	// handle responseStream
	for {
		message, error := responseStream.Recv()
		if error == io.EOF {
			// reached end of stream
			break
		}
		if error != nil {
			log.Fatalf("error while reading stream: %v", error)
		}

		log.Printf("Response from GreetManyTimes: %v", message.GetResult())
	}

	log.Printf("\t End of server streaming call \n\n")
}
