package main

import (
	"context"
	"fmt"
	"github.com/simplesteph/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
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

	firstName := request.Greeting.FirstName
	lastName := request.GetGreeting().LastName

	result := firstName + " " + lastName

	response := &greetpb.GreetResponse{Result: result}

	return response, nil
}

func (s server) GreetManyTimes(request *greetpb.GreetManyTimesRequest, timesServer greetpb.GreetService_GreetManyTimesServer) error {
	//TODO implement me
	panic("implement me")
}

func (s server) LongGreet(greetServer greetpb.GreetService_LongGreetServer) error {
	//TODO implement me
	panic("implement me")
}

func (s server) GreetEveryone(everyoneServer greetpb.GreetService_GreetEveryoneServer) error {
	//TODO implement me
	panic("implement me")
}

func (s server) GreetWithDeadline(context context.Context, request *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	//TODO implement me
	panic("implement me")
}
