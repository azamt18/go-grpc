package main

import (
	"context"
	"fmt"
	"github.com/simplesteph/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
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

	//makeUnaryCall(c)
	//makeServerStreamingCall(c)
	//makeClientStreamingCall(c)
	makeUnaryWithDeadline(c, 1000*time.Millisecond)
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

func makeClientStreamingCall(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do a client streaming RPC\n")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{FirstName: "Mike"}},
		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{FirstName: "John"}},
		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{FirstName: "Carl"}},
		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{FirstName: "Owen"}},
	}

	stream, error := c.LongGreet(context.Background())
	if error != nil {
		log.Fatalf("error while calling LongGreet: %v", error)
	}

	// iterate over slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	response, error := stream.CloseAndRecv()
	if error != nil {
		log.Fatalf("error while reading response from LongGreet: %v", error)
	}

	log.Printf("LongGreet response: %v", response.GetResult())

	fmt.Printf("Closing a client streaming RPC\n")
}

func makeBidirectionalStreamingCall(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting a BiDi Streaming RPC...\n")

	// create a stream by invoking the client
	stream, error := c.GreetEveryone(context.Background())
	if error != nil {
		log.Fatalf("error while creating stream: %v", error)
		return
	}

	// prepare a streaming request
	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{FirstName: "Mike"}},
		&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{FirstName: "John"}},
		&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{FirstName: "Carl"}},
		&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{FirstName: "Owen"}},
	}

	waitChannel := make(chan struct{})
	// send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		// close stream after sending messages
		stream.CloseSend()
	}()

	// receive a bunch of messages from the server (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading: %v", err)
				break
			}

			fmt.Printf("Response: %v\n", res.GetResult())
		}

		// when all is done, close waiting channel ->  "<-waitChannel" will be unblocked
		close(waitChannel)
	}()

	// block until everything is done
	<-waitChannel
}

func makeUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a UnaryWithDeadline RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Maarek",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("unexpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("error while calling GreetWithDeadline RPC: %v", err)
		}
		return
	}
	log.Printf("Response from GreetWithDeadline: %v", res.Result)
}
