package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"log"
	"net"
)

type server struct{}

func (s server) Sum(_ context.Context, request *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with: %v\n", request)

	firstNumber := request.GetFirstNumber()
	secondNumber := request.GetSecondNumber()

	result := firstNumber + secondNumber

	response := &calculatorpb.SumResponse{Result: result}

	return response, nil
}

func main() {
	fmt.Printf("start calc server\n")

	var listener, error = net.Listen("tcp", "0.0.0.0:50052")
	if error != nil {
		log.Fatalf("failed to listen %v", error)
	}

	grpcServer := grpc.NewServer()
	calculatorServiceServer := &server{}
	calculatorpb.RegisterCalculatorServiceServer(grpcServer, calculatorServiceServer)

	if error := grpcServer.Serve(listener); error != nil {
		log.Fatalf("failed to serve: %v\n", error)
	}
}
