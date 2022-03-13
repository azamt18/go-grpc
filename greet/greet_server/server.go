package main

import (
	"context"
	"fmt"
	"github.com/simplesteph/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
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
