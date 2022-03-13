package main

import (
	"context"
	"fmt"
	"github.com/simplesteph/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct{}

func main() {
	fmt.Println("start greet")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("FAILED TO LISTEN %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s server) Greet(context context.Context, request *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", request)

	firstName := request.GetGreeting().GetFirstName()
	lastName := request.GetGreeting().GetLastName()

	result := firstName + " " + lastName

	response := &greetpb.GreetResponse{Result: result}

	return response, nil
}

func (*server) GreetManyTimes(request *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with: %v\n", request)

	firstName := request.GetGreeting().GetFirstName()
	lastName := request.GetGreeting().GetLastName()

	for i := 0; i < 10; i++ {
		result := strconv.Itoa(i) + ": " + firstName + " " + lastName
		response := &greetpb.GreetManytimesResponse{Result: result}
		stream.SendMsg(response)
		//TODO: research: time changing does not affect
		time.Sleep(5 * time.Second)
	}

	return nil
}

func (s server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet was invoked with a streaming reqeuest\n")
	result := ""
	for {
		request, error := stream.Recv()
		if error == io.EOF {
			// finished reading the client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		// handle error
		if error != nil {
			log.Fatalf("error while reading client stream: %v\n", error)
		}

		firstName := request.GetGreeting().GetFirstName()
		lastName := request.GetGreeting().GetLastName()
		result += "Hello " + firstName + " " + lastName + "!\n"
	}
}

func (s server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with a streaming request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		firstName := req.Greeting.GetFirstName()
		lastName := req.Greeting.GetLastName()
		result := "Hello " + firstName + " " + lastName + "!\n"

		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("error while sending data to client: %v", sendErr)
			return err
		}
	}
}

func (s server) GreetWithDeadline(context context.Context, request *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	//TODO implement me
	panic("implement me")
}
